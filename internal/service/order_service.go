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
	cartRepo    repository.CartRepository
	db          *gorm.DB
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, cartRepo repository.CartRepository, db *gorm.DB) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		cartRepo:    cartRepo,
		db:          db,
	}
}

func (s *OrderService) Checkout(userID uint) (*models.Order, error) {
	cartItems, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil || len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	totalCents := 0
	var orderItems []models.OrderItem
	for _, item := range cartItems {
		product, err := s.productRepo.GetProductByIDForUpdate(tx, item.ProductID)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("product not found: " + item.Product.Name)
		}

		if product.Stock <= item.Quantity {
			tx.Rollback()
			return nil, errors.New("not enough stock for: " + item.Product.Name)
		}

		totalCents += product.Price * item.Quantity
		orderItems = append(orderItems, models.OrderItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}

	order := models.Order{
		UserID:     userID,
		TotalCents: totalCents,
		Status:     "paid",
		Items:      orderItems,
	}

	if err := s.orderRepo.CreateOrder(tx, &order); err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create order: " + err.Error())
	}

	if err := s.cartRepo.ClearCart(tx, userID); err != nil {
		tx.Rollback()
		return nil, errors.New("failed to clear cart: " + err.Error())
	}

	tx.Commit()
	return &order, nil
}
