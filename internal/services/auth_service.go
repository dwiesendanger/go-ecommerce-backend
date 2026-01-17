package services

import (
	"ecommerce-platform/internal/core/domain"
	"ecommerce-platform/pkg/auth"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) Register(email, password string) error {
	var existingUser domain.User
	if err := s.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return errors.New("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	user := domain.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.DB.Create(&user).Error
}

func (s *AuthService) Login(email, password string) (string, error) {
	var user domain.User

	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
