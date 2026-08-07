package purchases

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
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
}

type CreateInput struct {
	SupplierID uint        `json:"supplier_id" binding:"required"`
	Items      []ItemInput `json:"items" binding:"required,min=1"`
}

var ErrNoItems = errors.New("a purchase must include at least one item")

// Create records a new purchase and increases stock for every line item.
// The purchase, its items, and every resulting stock movement are written
// inside a single transaction so a failure midway leaves no partial state.
func (s *Service) Create(in CreateInput, userID uint) (*models.Purchase, error) {
	if len(in.Items) == 0 {
		return nil, ErrNoItems
	}

	var purchase models.Purchase
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var total float64
		for _, item := range in.Items {
			total += item.Quantity * item.Price
		}

		purchase = models.Purchase{
			ReferenceNumber: fmt.Sprintf("PO-%d", time.Now().UnixNano()/1000),
			SupplierID:      in.SupplierID,
			TotalAmount:     total,
			CreatedBy:       &userID,
		}
		if err := tx.Create(&purchase).Error; err != nil {
			return err
		}

		for _, item := range in.Items {
			purchaseItem := models.PurchaseItem{
				PurchaseID: purchase.ID,
				ProductID:  item.ProductID,
				Quantity:   item.Quantity,
				Price:      item.Price,
			}
			if err := tx.Create(&purchaseItem).Error; err != nil {
				return err
			}

			refID := purchase.ID
			if err := inventory.ApplyMovement(tx, item.ProductID, item.Quantity, models.MovementPurchase, "purchase", &refID,
				fmt.Sprintf("Purchase %s", purchase.ReferenceNumber), &userID); err != nil {
				return err
			}

			// Keep the product's cost basis current with the latest purchase price.
			if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
				Update("purchase_price", item.Price).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.Get(purchase.ID)
}

func (s *Service) List(limit int) ([]models.Purchase, error) {
	var purchases []models.Purchase
	q := s.db.Preload("Supplier").Preload("CreatedByUser.Role").Preload("Items").Preload("Items.Product").
		Order("created_at desc")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&purchases).Error
	return purchases, err
}

func (s *Service) Get(id uint) (*models.Purchase, error) {
	var purchase models.Purchase
	err := s.db.Preload("Supplier").Preload("CreatedByUser.Role").Preload("Items").Preload("Items.Product").
		First(&purchase, id).Error
	if err != nil {
		return nil, err
	}
	return &purchase, nil
}
