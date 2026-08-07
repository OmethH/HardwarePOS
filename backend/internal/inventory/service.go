// Package inventory is the single place allowed to change a product's
// current_stock. Every mutation goes through ApplyMovement so that a
// stock_movements row is always written alongside the stock change —
// per the product spec, stock must never be edited directly.
package inventory

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"hardwarepos/backend/internal/models"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// ApplyMovement adjusts a product's current_stock by quantity (positive to
// increase, negative to decrease) and records the movement. It must run
// inside the caller's transaction (tx) so stock changes and their triggering
// operation (sale, purchase, return) commit or roll back together.
func ApplyMovement(tx *gorm.DB, productID uint, quantity float64, movementType, referenceType string, referenceID *uint, note string, userID *uint) error {
	var product models.Product
	if err := tx.First(&product, productID).Error; err != nil {
		return fmt.Errorf("product %d not found: %w", productID, err)
	}

	newStock := product.CurrentStock + quantity
	if newStock < 0 {
		return fmt.Errorf("insufficient stock for %s: available %.2f, requested %.2f", product.Name, product.CurrentStock, -quantity)
	}

	if err := tx.Model(&models.Product{}).Where("id = ?", productID).
		Update("current_stock", newStock).Error; err != nil {
		return err
	}

	movement := models.StockMovement{
		ProductID:     productID,
		Quantity:      quantity,
		Type:          movementType,
		ReferenceType: referenceType,
		ReferenceID:   referenceID,
		Note:          note,
		CreatedBy:     userID,
	}
	return tx.Create(&movement).Error
}

func (s *Service) ListByProduct(productID uint) ([]models.StockMovement, error) {
	var movements []models.StockMovement
	err := s.db.Preload("CreatedByUser.Role").Where("product_id = ?", productID).
		Order("created_at desc").Find(&movements).Error
	return movements, err
}

func (s *Service) List(limit int) ([]models.StockMovement, error) {
	var movements []models.StockMovement
	q := s.db.Preload("Product").Preload("CreatedByUser.Role").Order("created_at desc")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&movements).Error
	return movements, err
}

var ErrInvalidAdjustment = errors.New("adjustment quantity cannot be zero")

// ManualAdjustment lets admins/managers correct stock counts (stocktake
// corrections, damage, etc). quantity is a signed delta.
func (s *Service) ManualAdjustment(productID uint, quantity float64, note string, userID uint) (*models.StockMovement, error) {
	if quantity == 0 {
		return nil, ErrInvalidAdjustment
	}

	var movement models.StockMovement
	err := s.db.Transaction(func(tx *gorm.DB) error {
		uid := userID
		if err := ApplyMovement(tx, productID, quantity, models.MovementAdjustment, "manual", nil, note, &uid); err != nil {
			return err
		}
		return tx.Where("product_id = ? AND type = ?", productID, models.MovementAdjustment).
			Order("created_at desc").First(&movement).Error
	})
	if err != nil {
		return nil, err
	}
	return &movement, nil
}
