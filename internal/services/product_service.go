package services

import (
	"context"
	"encoding/json"
	"time"

	"ecommerce-platform/internal/core/domain"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ProductService struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func NewProductService(db *gorm.DB, rdb *redis.Client) *ProductService {
	return &ProductService{DB: db, RDB: rdb}
}

// CreateProduct
func (s *ProductService) CreateProduct(input domain.Product) (*domain.Product, error) {
	if err := s.DB.Create(&input).Error; err != nil {
		return nil, err
	}
	go s.RDB.Del(context.Background(), "all_products")
	return &input, nil
}

// GetAllProducts
func (s *ProductService) GetAllProducts() ([]domain.Product, error) {
	ctx := context.Background()
	cacheKey := "all_products"

	val, err := s.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		var products []domain.Product
		if err := json.Unmarshal([]byte(val), &products); err == nil {
			return products, nil
		}
	}

	var products []domain.Product
	if err := s.DB.Find(&products).Error; err != nil {
		return nil, err
	}

	data, _ := json.Marshal(products)
	s.RDB.Set(ctx, cacheKey, data, 10*time.Minute)

	return products, nil
}
