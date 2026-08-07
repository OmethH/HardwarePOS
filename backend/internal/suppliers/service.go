package suppliers

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
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Company string `json:"company"`
}

func (s *Service) List(search string) ([]models.Supplier, error) {
	var suppliers []models.Supplier
	q := s.db.Order("name")
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("name LIKE ? OR company LIKE ? OR phone LIKE ?", like, like, like)
	}
	err := q.Find(&suppliers).Error
	return suppliers, err
}

func (s *Service) Get(id uint) (*models.Supplier, error) {
	var supplier models.Supplier
	if err := s.db.First(&supplier, id).Error; err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (s *Service) Create(in Input) (*models.Supplier, error) {
	supplier := models.Supplier{
		Name: in.Name, Phone: in.Phone, Email: in.Email, Address: in.Address, Company: in.Company,
	}
	if err := s.db.Create(&supplier).Error; err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (s *Service) Update(id uint, in Input) (*models.Supplier, error) {
	supplier, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	supplier.Name = in.Name
	supplier.Phone = in.Phone
	supplier.Email = in.Email
	supplier.Address = in.Address
	supplier.Company = in.Company
	if err := s.db.Save(supplier).Error; err != nil {
		return nil, err
	}
	return supplier, nil
}

func (s *Service) Delete(id uint) error {
	return s.db.Delete(&models.Supplier{}, id).Error
}
