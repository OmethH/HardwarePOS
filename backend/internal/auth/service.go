package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"hardwarepos/backend/internal/middleware"
	"hardwarepos/backend/internal/models"
)

var ErrInvalidCredentials = errors.New("invalid username or password")

type Service struct {
	db        *gorm.DB
	jwtSecret string
}

func NewService(db *gorm.DB, jwtSecret string) *Service {
	return &Service{db: db, jwtSecret: jwtSecret}
}

type LoginResult struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (s *Service) Login(username, password string) (*LoginResult, error) {
	var user models.User
	if err := s.db.Preload("Role").Where("username = ? AND active = 1", username).First(&user).Error; err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.issueToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginResult{Token: token, User: user}, nil
}

func (s *Service) issueToken(user models.User) (string, error) {
	claims := middleware.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *Service) Me(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Role").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}
