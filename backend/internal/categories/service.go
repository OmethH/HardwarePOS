package categories

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

func (s *Service) List() ([]models.Category, error) {
	var categories []models.Category
	err := s.db.Order("name").Find(&categories).Error
	return categories, err
}

func (s *Service) Get(id uint) (*models.Category, error) {
	var category models.Category
	if err := s.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *Service) Create(name string) (*models.Category, error) {
	category := models.Category{Name: name}
	if err := s.db.Create(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *Service) Update(id uint, name string) (*models.Category, error) {
	category, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	category.Name = name
	if err := s.db.Save(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (s *Service) Delete(id uint) error {
	return s.db.Delete(&models.Category{}, id).Error
}
