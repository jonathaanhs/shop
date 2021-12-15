package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/learn/shop/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewOrderRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbx := sqlx.NewDb(db, "sqlmock")

	orderRepo := repository.NewOrderRepository(dbx)

	assert.NotNil(t, orderRepo)
}
