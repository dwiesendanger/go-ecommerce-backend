package services

import (
	"ecommerce-platform/internal/core/domain"
	"ecommerce-platform/internal/workers"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CheckoutService struct {
	DB        *gorm.DB
	OrderChan chan<- workers.OrderJob // "Send-only" Channel
}

func NewCheckoutService(db *gorm.DB, orderChan chan<- workers.OrderJob) *CheckoutService {
	return &CheckoutService{
		DB:        db,
		OrderChan: orderChan,
	}
}

func (s *CheckoutService) PlaceOrder(userID uint) (*domain.Order, error) {
	var order domain.Order
	var userEmail string

	err := s.DB.Transaction(func(tx *gorm.DB) error {

		var cart domain.Cart
		if err := tx.Preload("Items.Product").First(&cart, "user_id = ?", userID).Error; err != nil {
			return errors.New("cart not found or empty")
		}
		if len(cart.Items) == 0 {
			return errors.New("cart is empty")
		}

		var user domain.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}
		userEmail = user.Email

		var totalAmount float64
		var orderItems []domain.OrderItem

		for _, item := range cart.Items {
			if item.Product.Stock < item.Quantity {
				return fmt.Errorf("not enough stock for product: %s", item.Product.Name)
			}

			newStock := item.Product.Stock - item.Quantity
			if err := tx.Model(&domain.Product{}).Where("id = ?", item.Product.ID).Update("stock", newStock).Error; err != nil {
				return err
			}

			totalAmount += item.Product.Price * float64(item.Quantity)

			orderItems = append(orderItems, domain.OrderItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.Product.Price,
			})
		}

		order = domain.Order{
			UserID:      userID,
			OrderNumber: fmt.Sprintf("ORD-%d", time.Now().UnixNano()),
			Status:      "pending",
			TotalAmount: totalAmount,
			Items:       orderItems,
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		if err := tx.Where("cart_id = ?", cart.ID).Delete(&domain.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.OrderChan <- workers.OrderJob{
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		UserEmail:   userEmail,
	}

	return &order, nil
}

func (s *CheckoutService) GetOrdersByUser(userID uint) ([]domain.Order, error) {
	var orders []domain.Order

	err := s.DB.Preload("Items.Product").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&orders).Error

	if err != nil {
		return nil, err
	}
	return orders, nil
}
