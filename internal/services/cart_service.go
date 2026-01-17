package services

import (
	"ecommerce-platform/internal/core/domain"
	"errors"

	"gorm.io/gorm"
)

type CartService struct {
	DB *gorm.DB
}

func NewCartService(db *gorm.DB) *CartService {
	return &CartService{DB: db}
}

func (s *CartService) AddItem(userID uint, productID uint, quantity int) error {
	var cart domain.Cart

	if err := s.DB.FirstOrCreate(&cart, domain.Cart{UserID: userID}).Error; err != nil {
		return err
	}

	var product domain.Product
	if err := s.DB.First(&product, productID).Error; err != nil {
		return errors.New("product not found")
	}

	var cartItem domain.CartItem
	result := s.DB.Where("cart_id = ? AND product_id = ?", cart.ID, productID).First(&cartItem)

	if result.Error == nil {
		cartItem.Quantity += quantity
		return s.DB.Save(&cartItem).Error
	}

	newItem := domain.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  quantity,
	}
	return s.DB.Create(&newItem).Error
}

func (s *CartService) GetCart(userID uint) (*domain.Cart, error) {
	var cart domain.Cart

	err := s.DB.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error

	if err != nil {
		return nil, errors.New("cart empty")
	}
	return &cart, nil
}
