package repository

import (
	"game-store-api/internal/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	AddItem(item *models.CartItem) error
	GetCartByUserID(userID uint) ([]models.CartItem, error)
	RemoveItem(userID, productID uint) error
	ClearCart(tx *gorm.DB, userID uint) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) AddItem(item *models.CartItem) error {
	// If item exists update quantity else create
	var existingItem models.CartItem
	err := r.db.Where("user_id = ? AND product_id = ?", item.UserID, item.ProductID).First(&existingItem).Error
	if err == nil {
		existingItem.Quantity += item.Quantity
		if existingItem.Quantity <= 0 {
			return r.db.Delete(&existingItem).Error
		}
		return r.db.Save(&existingItem).Error
	}
	if item.Quantity <= 0 {
		return nil
	}
	return r.db.Create(item).Error
}

func (r *cartRepository) GetCartByUserID(userID uint) ([]models.CartItem, error) {
	var CartItems []models.CartItem
	err := r.db.Preload("Product").Where("user_id = ?", userID).Find(&CartItems).Error
	return CartItems, err
}

func (r *cartRepository) RemoveItem(userID, productID uint) error {
	var existingItem models.CartItem
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&existingItem).Error
	if err != nil {
		return err
	}
	return r.db.Delete(&existingItem).Error
}

func (r *cartRepository) ClearCart(tx *gorm.DB, userID uint) error {
	return tx.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error
}
