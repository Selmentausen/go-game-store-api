package service

import (
	"game-store-api/internal/models"
	"game-store-api/internal/repository"
)

type ProductService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s ProductService) CreateProduct(name, description, sku string, price, stock int) error {
	product := models.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		SKU:         sku,
	}
	return s.productRepo.CreateProduct(&product)
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.productRepo.GetAllProducts()
}
