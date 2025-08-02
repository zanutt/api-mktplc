package product

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

func ListProducts(db *gorm.DB, nameFilter, categoryFilter string, page, limit int) ([]Product, int64, error) {
	var products []Product
	query := db.Model(&Product{})

	if nameFilter != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(nameFilter)+"%")
	}

	if categoryFilter != "" {
		query = query.Where("LOWER(category) = ?", strings.ToLower(categoryFilter))
	}

	var totalCount int64
	query.Count(&totalCount)

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, totalCount, nil
}

func UpdateProduct(db *gorm.DB, updated Product) (Product, error) {
	if updated.ID == 0 || updated.Name == "" || updated.Price <= 0 {
		return Product{}, errors.New("invalid product data")
	}

	var existing Product
	if err := db.First(&existing, updated.ID).Error; err != nil {
		return Product{}, err
	}

	existing.Name = updated.Name
	existing.Desc = updated.Desc
	existing.Price = updated.Price
	existing.Category = updated.Category

	if err := db.Save(&existing).Error; err != nil {
		return Product{}, err
	}

	return existing, nil
}

func DeleteProduct(db *gorm.DB, id uint) error {
	if id == 0 {
		return errors.New("invalid ID")
	}

	if err := db.Delete(&Product{}, id).Error; err != nil {
		return err
	}

	return nil
}
