package returns

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"hardwarepos/backend/internal/inventory"
	"hardwarepos/backend/internal/models"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type ItemInput struct {
	SaleItemID uint    `json:"sale_item_id" binding:"required"`
	Quantity   float64 `json:"quantity" binding:"required"`
}

type CreateInput struct {
	SaleID uint        `json:"sale_id" binding:"required"`
	Items  []ItemInput `json:"items" binding:"required,min=1"`
}

var (
	ErrNoItems          = errors.New("a return must include at least one item")
	ErrExceedsPurchased = errors.New("return quantity exceeds the quantity sold")
)

// Create processes a return: validates the requested quantities against
// what's left to return on the original sale, increases stock for each
// item, and marks the sale COMPLETED / PARTIALLY_RETURNED / RETURNED.
// Sales are never deleted or edited in place — this only ever adds records.
func (s *Service) Create(in CreateInput, userID uint) (*models.Return, error) {
	if len(in.Items) == 0 {
		return nil, ErrNoItems
	}

	var ret models.Return
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var sale models.Sale
		if err := tx.Preload("Items").First(&sale, in.SaleID).Error; err != nil {
			return fmt.Errorf("sale not found: %w", err)
		}

		saleItemByID := make(map[uint]models.SaleItem, len(sale.Items))
		for _, si := range sale.Items {
			saleItemByID[si.ID] = si
		}

		var totalRefund float64
		ret = models.Return{
			ReturnNumber: fmt.Sprintf("RET-%d", time.Now().UnixNano()/1000),
			SaleID:       sale.ID,
			CreatedBy:    &userID,
		}
		if err := tx.Create(&ret).Error; err != nil {
			return err
		}

		for _, item := range in.Items {
			saleItem, ok := saleItemByID[item.SaleItemID]
			if !ok || saleItem.SaleID != sale.ID {
				return fmt.Errorf("sale item %d does not belong to sale %d", item.SaleItemID, sale.ID)
			}
			remaining := saleItem.Quantity - saleItem.ReturnedQuantity
			if item.Quantity > remaining {
				return fmt.Errorf("%w: item %d (remaining %.2f, requested %.2f)", ErrExceedsPurchased, item.SaleItemID, remaining, item.Quantity)
			}

			refund := item.Quantity * saleItem.SellingPrice
			totalRefund += refund

			returnItem := models.ReturnItem{
				ReturnID:     ret.ID,
				SaleItemID:   saleItem.ID,
				ProductID:    saleItem.ProductID,
				Quantity:     item.Quantity,
				RefundAmount: refund,
			}
			if err := tx.Create(&returnItem).Error; err != nil {
				return err
			}

			if err := tx.Model(&models.SaleItem{}).Where("id = ?", saleItem.ID).
				Update("returned_quantity", saleItem.ReturnedQuantity+item.Quantity).Error; err != nil {
				return err
			}

			refID := ret.ID
			if err := inventory.ApplyMovement(tx, saleItem.ProductID, item.Quantity, models.MovementReturn, "return", &refID,
				fmt.Sprintf("Return %s", ret.ReturnNumber), &userID); err != nil {
				return err
			}

			saleItemByID[item.SaleItemID] = models.SaleItem{
				ID: saleItem.ID, Quantity: saleItem.Quantity, ReturnedQuantity: saleItem.ReturnedQuantity + item.Quantity,
			}
		}

		if err := tx.Model(&ret).Update("total_refund", totalRefund).Error; err != nil {
			return err
		}

		newStatus := models.SaleStatusPartiallyReturned
		fullyReturned := true
		for _, si := range saleItemByID {
			if si.ReturnedQuantity < si.Quantity {
				fullyReturned = false
				break
			}
		}
		if fullyReturned {
			newStatus = models.SaleStatusReturned
		}
		if err := tx.Model(&sale).Update("status", newStatus).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.Get(ret.ID)
}

func (s *Service) List(limit int) ([]models.Return, error) {
	var returns []models.Return
	q := s.db.Preload("Sale").Preload("CreatedByUser.Role").Preload("Items").Preload("Items.Product").Order("created_at desc")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&returns).Error
	return returns, err
}

func (s *Service) Get(id uint) (*models.Return, error) {
	var ret models.Return
	err := s.db.Preload("Sale").Preload("CreatedByUser.Role").Preload("Items").Preload("Items.Product").First(&ret, id).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
