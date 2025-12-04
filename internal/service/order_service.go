package service

import (
	"errors"
	"game-store-api/internal/models"
	"game-store-api/internal/repository"

	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	db          *gorm.DB
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, db *gorm.DB) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		db:          db,
	}
}

func (s *OrderService) BuyProduct(userID, productID uint) (*models.Order, error) {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	product, err := s.productRepo.GetProductByIDForUpdate(tx, productID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("product not found")
	}

	if product.Stock <= 0 {
		tx.Rollback()
		return nil, errors.New("product out of stock")
	}

	order := models.Order{
		UserID:    userID,
		ProductID: productID,
	}
	if err := s.orderRepo.CreateOrder(tx, &order); err != nil {
		tx.Rollback()
		return nil, err
	}

	product.Stock = product.Stock - 1
	if err := s.productRepo.UpdateProduct(tx, product); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &order, nil
}
