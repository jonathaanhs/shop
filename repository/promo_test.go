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

func TestNewPromoRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbx := sqlx.NewDb(db, "sqlmock")

	promoRepo := repository.NewPromoRepository(dbx)

	assert.NotNil(t, promoRepo)
}

func TestGetPromoByProductID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbx := sqlx.NewDb(db, "sqlmock")

	promoRepo := repository.NewPromoRepository(dbx)

	query := regexp.QuoteMeta("select promo_id, product_id, promo_type, reward, min_qty from promos where product_id = $1")
	productID := rand.Int63()

	var response repository.Promo
	err = faker.FakeData(&response)
	if err != nil {
		t.Fatalf("error while fake data %+v", err)
	}

	rowsResponse := sqlmock.NewRows([]string{"promo_id", "product_id", "promo_type", "reward", "min_qty"}).AddRow(response.PromoID, response.ProductID, response.PromoType, response.Reward, response.MinQty)
	mock.ExpectQuery(query).WithArgs(productID).WillReturnRows(rowsResponse).RowsWillBeClosed()

	output, err := promoRepo.GetPromoByProductID(context.Background(), productID)
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}

func TestGetAllPromo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbx := sqlx.NewDb(db, "sqlmock")

	promoRepo := repository.NewPromoRepository(dbx)

	query := regexp.QuoteMeta("select promo_id, product_id, promo_type, reward, min_qty from promos order by promo_id asc")

	var response repository.Promo
	err = faker.FakeData(&response)
	if err != nil {
		t.Fatalf("error while fake data %+v", err)
	}

	rowsResponse := sqlmock.NewRows([]string{"promo_id", "product_id", "promo_type", "reward", "min_qty"}).AddRow(response.PromoID, response.ProductID, response.PromoType, response.Reward, response.MinQty)
	mock.ExpectQuery(query).WillReturnRows(rowsResponse).RowsWillBeClosed()

	output, err := promoRepo.GetAllPromo(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}
