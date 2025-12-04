package repository

import (
	"game-store-api/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) error
	GetAllProducts() ([]models.Product, error)
	GetProductByID(id uint) (*models.Product, error)
	GetProductByIDForUpdate(tx *gorm.DB, id uint) (*models.Product, error)
	UpdateProduct(tx *gorm.DB, product *models.Product) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *productRepository) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	return &product, err
}

func (r *productRepository) GetProductByIDForUpdate(tx *gorm.DB, id uint) (*models.Product, error) {
	var product models.Product
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, id).Error
	return &product, err
}

func (r *productRepository) UpdateProduct(tx *gorm.DB, product *models.Product) error {
	return tx.Save(product).Error
}
