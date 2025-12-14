package service

import (
	"context"
	"errors"
	pb "game-store-api/internal/grpc/payment"
	"game-store-api/internal/models"
	"game-store-api/internal/repository"
	"time"

	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo     repository.OrderRepository
	productRepo   repository.ProductRepository
	cartRepo      repository.CartRepository
	paymentClient pb.PaymentServiceClient
	db            *gorm.DB
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	cartRepo repository.CartRepository,
	paymentClient pb.PaymentServiceClient,
	db *gorm.DB) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		productRepo:   productRepo,
		cartRepo:      cartRepo,
		paymentClient: paymentClient,
		db:            db,
	}
}

func (s *OrderService) Checkout(userID uint) (*models.Order, error) {
	cartItems, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil || len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// --- Make a payment request ---
	totalCents := 0
	for _, item := range cartItems {
		totalCents += item.Quantity * item.Product.Price
	}

	paymentReq := &pb.PaymentRequest{
		OrderId:        int64(userID),
		Amount:         float32(totalCents) / 100.0,
		Currency:       "USD",
		CredCardNumber: "1212-1212-1212-1212",
	}

	paymentRes, err := s.paymentClient.ProcessPayment(ctx, paymentReq)
	if err != nil {
		return nil, errors.New("payment service unavailable")
	}

	if !paymentRes.Success {
		return nil, errors.New("payment declined: " + paymentRes.Message)
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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
		product.Stock -= item.Quantity

		if err := s.productRepo.UpdateProduct(tx, product); err != nil {
			tx.Rollback()
			return nil, errors.New("product update failed: " + item.Product.Name)
		}

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
