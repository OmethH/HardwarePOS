package users

import (
	"gorm.io/gorm"

	"hardwarepos/backend/internal/auth"
	"hardwarepos/backend/internal/models"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type CreateInput struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

type UpdateInput struct {
	Name     string `json:"name" binding:"required"`
	RoleID   uint   `json:"role_id" binding:"required"`
	Active   bool   `json:"active"`
	Password string `json:"password"` // optional: reset password when non-empty
}

func (s *Service) List() ([]models.User, error) {
	var users []models.User
	err := s.db.Preload("Role").Order("name").Find(&users).Error
	return users, err
}

func (s *Service) Get(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) Create(in CreateInput) (*models.User, error) {
	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}
	user := models.User{
		Name:         in.Name,
		Username:     in.Username,
		PasswordHash: hash,
		RoleID:       in.RoleID,
		Active:       true,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return s.Get(user.ID)
}

func (s *Service) Update(id uint, in UpdateInput) (*models.User, error) {
	user, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	user.Name = in.Name
	user.RoleID = in.RoleID
	user.Active = in.Active
	if in.Password != "" {
		hash, err := auth.HashPassword(in.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hash
	}
	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}
	return s.Get(user.ID)
}

// Deactivate is used instead of a hard delete so historical sales/purchases
// still resolve their created_by user.
func (s *Service) Deactivate(id uint) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("active", false).Error
}

func (s *Service) ListRoles() ([]models.Role, error) {
	var roles []models.Role
	err := s.db.Find(&roles).Error
	return roles, err
}
