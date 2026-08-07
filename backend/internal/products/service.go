package products

import (
	"errors"

	"gorm.io/gorm"

	"hardwarepos/backend/internal/models"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type Input struct {
	SKU           string  `json:"sku" binding:"required"`
	Barcode       *string `json:"barcode"`
	Name          string  `json:"name" binding:"required"`
	CategoryID    *uint   `json:"category_id"`
	Brand         string  `json:"brand"`
	Description   string  `json:"description"`
	UnitType      string  `json:"unit_type" binding:"required"`
	PurchasePrice float64 `json:"purchase_price"`
	SellingPrice  float64 `json:"selling_price"`
	MinimumStock  float64 `json:"minimum_stock"`
	Status        string  `json:"status"`
}

type ListFilter struct {
	Search     string
	CategoryID uint
	LowStock   bool
	Status     string
}

func (s *Service) List(f ListFilter) ([]models.Product, error) {
	var products []models.Product
	q := s.db.Preload("Category").Order("name")

	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("name LIKE ? OR sku LIKE ? OR barcode LIKE ?", like, like, like)
	}
	if f.CategoryID != 0 {
		q = q.Where("category_id = ?", f.CategoryID)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.LowStock {
		q = q.Where("current_stock <= minimum_stock")
	}

	err := q.Find(&products).Error
	return products, err
}

func (s *Service) Get(id uint) (*models.Product, error) {
	var product models.Product
	if err := s.db.Preload("Category").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *Service) GetByBarcode(barcode string) (*models.Product, error) {
	var product models.Product
	if err := s.db.Preload("Category").Where("barcode = ?", barcode).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *Service) Create(in Input) (*models.Product, error) {
	if in.UnitType == "" {
		in.UnitType = "Piece"
	}
	status := in.Status
	if status == "" {
		status = "active"
	}
	product := models.Product{
		SKU:           in.SKU,
		Barcode:       in.Barcode,
		Name:          in.Name,
		CategoryID:    in.CategoryID,
		Brand:         in.Brand,
		Description:   in.Description,
		UnitType:      in.UnitType,
		PurchasePrice: in.PurchasePrice,
		SellingPrice:  in.SellingPrice,
		MinimumStock:  in.MinimumStock,
		CurrentStock:  0, // stock only ever changes via inventory.ApplyMovement
		Status:        status,
	}
	if err := s.db.Create(&product).Error; err != nil {
		return nil, err
	}
	return s.Get(product.ID)
}

var ErrPurchasePriceLocked = errors.New("only admins and managers may change the purchase price")

// Update edits product metadata. current_stock is intentionally excluded —
// it can only move through inventory.ApplyMovement (purchases, sales,
// returns, manual adjustments).
func (s *Service) Update(id uint, in Input, allowPriceChange bool) (*models.Product, error) {
	product, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if !allowPriceChange && in.PurchasePrice != product.PurchasePrice {
		return nil, ErrPurchasePriceLocked
	}

	product.SKU = in.SKU
	product.Barcode = in.Barcode
	product.Name = in.Name
	product.CategoryID = in.CategoryID
	product.Brand = in.Brand
	product.Description = in.Description
	product.UnitType = in.UnitType
	product.SellingPrice = in.SellingPrice
	product.MinimumStock = in.MinimumStock
	if in.Status != "" {
		product.Status = in.Status
	}
	if allowPriceChange {
		product.PurchasePrice = in.PurchasePrice
	}

	if err := s.db.Save(product).Error; err != nil {
		return nil, err
	}
	return s.Get(product.ID)
}

func (s *Service) Delete(id uint) error {
	return s.db.Delete(&models.Product{}, id).Error
}

func (s *Service) LowStockCount() (int64, error) {
	var count int64
	err := s.db.Model(&models.Product{}).Where("current_stock <= minimum_stock AND current_stock > 0").Count(&count).Error
	return count, err
}

func (s *Service) OutOfStockCount() (int64, error) {
	var count int64
	err := s.db.Model(&models.Product{}).Where("current_stock <= 0").Count(&count).Error
	return count, err
}

func (s *Service) TotalCount() (int64, error) {
	var count int64
	err := s.db.Model(&models.Product{}).Count(&count).Error
	return count, err
}
