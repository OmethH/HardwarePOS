package settings

import (
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
	ShopName      string  `json:"shop_name" binding:"required"`
	Address       string  `json:"address"`
	Phone         string  `json:"phone"`
	ReceiptFooter string  `json:"receipt_footer"`
	TaxPercentage float64 `json:"tax_percentage"`
}

func (s *Service) Get() (*models.Settings, error) {
	var settings models.Settings
	if err := s.db.First(&settings, 1).Error; err != nil {
		return nil, err
	}
	return &settings, nil
}

func (s *Service) Update(in Input) (*models.Settings, error) {
	settings, err := s.Get()
	if err != nil {
		return nil, err
	}
	settings.ShopName = in.ShopName
	settings.Address = in.Address
	settings.Phone = in.Phone
	settings.ReceiptFooter = in.ReceiptFooter
	settings.TaxPercentage = in.TaxPercentage
	if err := s.db.Save(settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}
