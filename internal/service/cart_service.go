package service

import (
	"game-store-api/internal/models"
	"game-store-api/internal/repository"
)

type CartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) *CartService {
	return &CartService{cartRepo: cartRepo, productRepo: productRepo}
}

func (s *CartService) AddToCart(userID, productID uint, quantity int) error {
	_, err := s.productRepo.GetProductByID(productID)
	if err != nil {
		return err
	}

	item := models.CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}
	return s.cartRepo.AddItem(&item)
}

func (s *CartService) GetCart(userID uint) ([]models.CartItem, error) {
	return s.cartRepo.GetCartByUserID(userID)
}

func (s *CartService) RemoveItem(userID, productID uint) error {
	return s.cartRepo.RemoveItem(userID, productID)
}
