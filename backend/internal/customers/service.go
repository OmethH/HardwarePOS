package customers

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
	Address string `json:"address"`
}

func (s *Service) List(search string) ([]models.Customer, error) {
	var customers []models.Customer
	q := s.db.Order("name")
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("name LIKE ? OR phone LIKE ?", like, like)
	}
	err := q.Find(&customers).Error
	return customers, err
}

func (s *Service) Get(id uint) (*models.Customer, error) {
	var customer models.Customer
	if err := s.db.First(&customer, id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (s *Service) Create(in Input) (*models.Customer, error) {
	customer := models.Customer{Name: in.Name, Phone: in.Phone, Address: in.Address}
	if err := s.db.Create(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (s *Service) Update(id uint, in Input) (*models.Customer, error) {
	customer, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	customer.Name = in.Name
	customer.Phone = in.Phone
	customer.Address = in.Address
	if err := s.db.Save(customer).Error; err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *Service) Delete(id uint) error {
	if id == 1 {
		return gorm.ErrInvalidData // protect the default "Walk-in Customer"
	}
	return s.db.Delete(&models.Customer{}, id).Error
}
