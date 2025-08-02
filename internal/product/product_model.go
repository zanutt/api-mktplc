package product

import "errors"

type Product struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Desc     string `gorm:"column:description"`
	Price    float64
	Category string
}

func NewProduct(name, desc string, price float64, category string) (Product, error) {
	if name == "" || price <= 0 {
		return Product{}, errors.New("invalid product data")
	}
	return Product{
		Name:     name,
		Desc:     desc,
		Price:    price,
		Category: category,
	}, nil
}
