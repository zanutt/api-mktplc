package product

import (
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

func TestNewProduct_SaveToDB(t *testing.T) {
	p, err := NewProduct("Teclado", "TecladoRGB", 299.99, "Periféricos")
	assert.Nil(t, err)

	result := testDB.Create(&p)
	assert.Nil(t, result.Error)

	var fetched Product
	testDB.First(&fetched, "name = ?", "Teclado")

	assert.Equal(t, p.Desc, fetched.Desc)

}

func TestListProducts(t *testing.T) {
	// Reseta o db
	testDB.Exec("DELETE FROM products")

	// Cria products de exemplo
	testDB.Create(&Product{Name: "Notebook", Desc: "i7, 16GB RAM", Price: 4500.00, Category: "Eletrônicos"})
	testDB.Create(&Product{Name: "Cafeteira", Desc: "Expresso", Price: 320.00, Category: "Eletrodomésticos"})
	testDB.Create(&Product{Name: "Mouse Gamer", Desc: "RGB", Price: 159.90, Category: "Periféricos"})

	products, err := ListProducts(testDB, "", "")
	assert.Nil(t, err)
	assert.Len(t, products, 3)
}

func TestListProducts_WithFilters(t *testing.T) {
	testDB.Exec("DELETE FROM products")

	testDB.Create(&Product{Name: "Notebook Gamer", Desc: "", Price: 6000.00, Category: "Eletrônicos"})
	testDB.Create(&Product{Name: "Notebook Office", Desc: "", Price: 3200.00, Category: "Eletrônicos"})
	testDB.Create(&Product{Name: "Furadeira", Desc: "", Price: 250.00, Category: "Ferramentas"})

	// busca por nome que contenha "note" e categoria "Eletrônicos"
	products, err := ListProducts(testDB, "note", "Eletrônicos")
	assert.Nil(t, err)
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
