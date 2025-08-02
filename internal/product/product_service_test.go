package product

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewProduct_Valid(t *testing.T) {
	p, err := NewProduct("Gaming Mouse", "8000 DPI Sensor", 159.90, "Peripherals")

	assert.Nil(t, err)
	assert.Equal(t, "Gaming Mouse", p.Name)
}

func TestNewProduct_Invalid(t *testing.T) {
	_, err := NewProduct("", "Some desc", 0, "Peripherals")

	assert.NotNil(t, err)
	assert.Equal(t, "invalid product data", err.Error())
}

func TestListProducts(t *testing.T) {
	// Reseta o db
	testDB.Exec("DELETE FROM products")

	// Cria produtos de exemplo
	testDB.Create(&Product{Name: "Notebook", Desc: "i7, 16GB RAM", Price: 4500.00, Category: "Eletrônicos"})
	testDB.Create(&Product{Name: "Cafeteira", Desc: "Expresso", Price: 320.00, Category: "Eletrodomésticos"})
	testDB.Create(&Product{Name: "Mouse Gamer", Desc: "RGB", Price: 159.90, Category: "Periféricos"})

	products, total, err := ListProducts(testDB, "", "", 1, 10)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, products, 3)
}

func TestListProducts_WithFilters(t *testing.T) {
	testDB.Exec("DELETE FROM products")

	testDB.Create(&Product{Name: "Notebook Gamer", Desc: "", Price: 6000.00, Category: "Eletrônicos"})
	testDB.Create(&Product{Name: "Notebook Office", Desc: "", Price: 3200.00, Category: "Eletrônicos"})
	testDB.Create(&Product{Name: "Furadeira", Desc: "", Price: 250.00, Category: "Ferramentas"})

	products, total, err := ListProducts(testDB, "note", "Eletrônicos", 1, 10)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, products, 2)
}

func TestUpdateProduct(t *testing.T) {
	testDB.Exec("DELETE FROM products")

	// cria produto
	original := Product{
		Name:     "Monitor Full HD",
		Desc:     "60Hz",
		Price:    799.90,
		Category: "Periféricos",
	}
	testDB.Create(&original)

	// prepara atualização
	update := Product{
		ID:       original.ID,
		Name:     "Monitor Full HD 75Hz",
		Desc:     "Atualizado para 75Hz",
		Price:    899.90,
		Category: "Periféricos",
	}

	updated, err := UpdateProduct(testDB, update)

	assert.Nil(t, err)
	assert.Equal(t, update.Name, updated.Name)
	assert.Equal(t, update.Desc, updated.Desc)
	assert.Equal(t, update.Price, updated.Price)
}

func TestDeleteProduct(t *testing.T) {
	testDB.Exec("DELETE FROM products")

	// cria o product
	product := Product{
		Name:     "Webcam HD",
		Desc:     "720p",
		Price:    249.90,
		Category: "Periféricos",
	}
	testDB.Create(&product)

	// deleta produto
	err := DeleteProduct(testDB, product.ID)
	assert.Nil(t, err)

	// confere se o produto esta ainda
	var p Product
	result := testDB.First(&p, product.ID)
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestListProducts_WithPagination(t *testing.T) {
	testDB.Exec("DELETE FROM products")

	// Insere 6 produtos
	for i := 1; i <= 6; i++ {
		testDB.Create(&Product{
			Name:     fmt.Sprintf("Produto %d", i),
			Desc:     "Teste",
			Price:    100,
			Category: "Geral",
		})
	}

	page1, total1, err := ListProducts(testDB, "", "", 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, int64(6), total1)
	assert.Len(t, page1, 2)

	page2, total2, err := ListProducts(testDB, "", "", 2, 2)
	assert.Nil(t, err)
	assert.Equal(t, int64(6), total2)
	assert.Len(t, page2, 2)

	page3, total3, err := ListProducts(testDB, "", "", 3, 2)
	assert.Nil(t, err)
	assert.Equal(t, int64(6), total3)
	assert.Len(t, page3, 2)

	page4, total4, err := ListProducts(testDB, "", "", 4, 2)
	assert.Nil(t, err)
	assert.Equal(t, int64(6), total4)
	assert.Len(t, page4, 0)
}
