package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zanutt/api-mktplc/internal/product"
)

func TestListProducts(t *testing.T) {
	SetupTest(t)
	// Cria produto de exemplo
	TestDB.Create(&product.Product{
		Name:     "Notebook",
		Desc:     "i7, 16GB RAM",
		Price:    4500.00,
		Category: "Eletrônicos",
	})

	products, total, err := product.ListProducts(TestDB, "", "", 1, 10)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, products, 1)
}
