package repository_test

import (
	"context"
	"math/rand"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker"
	"github.com/jmoiron/sqlx"
	"github.com/learn/shop/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewProductRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbx := sqlx.NewDb(db, "sqlmock")

	productRepo := repository.NewProductRepository(dbx)

	assert.NotNil(t, productRepo)
}

func TestGetProductByProductID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbx := sqlx.NewDb(db, "sqlmock")

	productRepo := repository.NewProductRepository(dbx)

	query := regexp.QuoteMeta("select product_id, sku, name, price, qty from products where product_id = $1")
	productID := rand.Int63()

	var response repository.Product
	err = faker.FakeData(&response)
	if err != nil {
		t.Fatalf("error while fake data %+v", err)
	}

	rowsResponse := sqlmock.NewRows([]string{"product_id", "sku", "name", "price", "qty"}).AddRow(response.ProductID, response.Sku, response.Name, response.Price, response.Qty)
	mock.ExpectQuery(query).WithArgs(productID).WillReturnRows(rowsResponse).RowsWillBeClosed()

	output, err := productRepo.GetProductByProductID(context.Background(), productID)
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}

func TestGetAllProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbx := sqlx.NewDb(db, "sqlmock")

	productRepo := repository.NewProductRepository(dbx)

	query := regexp.QuoteMeta("select product_id, sku, name, price, qty from products order by product_id asc")

	var response repository.Product
	err = faker.FakeData(&response)
	if err != nil {
		t.Fatalf("error while fake data %+v", err)
	}

	rowsResponse := sqlmock.NewRows([]string{"product_id", "sku", "name", "price", "qty"}).AddRow(response.ProductID, response.Sku, response.Name, response.Price, response.Qty)
	mock.ExpectQuery(query).WillReturnRows(rowsResponse).RowsWillBeClosed()

	output, err := productRepo.GetAllProduct(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}
