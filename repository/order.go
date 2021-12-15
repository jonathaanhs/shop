package repository

import (
	"context"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/learn/shop/helper"
)

type orderRepo struct {
	DB *sqlx.DB
}

type Order struct {
	OrderID int64     `json:"order_id" db:"order_id"`
	Date    time.Time `json:"date" db:"date"`
	Total   float64   `json:"total" db:"total"`
}

type OrderDetail struct {
	OrderDetailID int64   `json:"order_detail_id" db:"order_detail_id"`
	OrderID       int64   `json:"order_id" db:"order_id"`
	ProductID     int64   `json:"product_id" db:"product_id"`
	PromoID       int64   `json:"promo_id" db:"promo_id"`
	Price         float64 `json:"price" db:"price"`
	Qty           int64   `json:"qty" db:"qty"`
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return orderRepo{
		DB: db,
	}
}

func (o orderRepo) CreateOrder(tx *sqlx.Tx, ctx context.Context, form Order) (orderID int64, err error) {
	err = tx.QueryRowxContext(ctx, "insert into orders(date, total) values($1, $2) RETURNING order_id", form.Date, form.Total).Scan(&orderID)
	if err != nil {
		return orderID, err
	}

	return orderID, nil
}

func (o orderRepo) CreateOrderDetails(tx *sqlx.Tx, ctx context.Context, form []OrderDetail) (err error) {
	sqlInsert := "insert into order_details(order_id, product_id, promo_id, price, qty) values"
	rowSQL := "(?, ?, ?, ?, ?)"

	vals := []interface{}{}
	var inserts []string

	for _, val := range form {
		vals = append(vals, val.OrderID, val.ProductID, val.PromoID, val.Price, val.Qty)
		inserts = append(inserts, rowSQL)
	}

	sqlInsert = sqlInsert + strings.Join(inserts, ",")
	sqlInsert = helper.ReplaceSQL(sqlInsert, "?")

	stmt, err := tx.PrepareContext(ctx, sqlInsert)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, vals...)
	if err != nil {
		return err
	}

	return nil
}

func (o orderRepo) BeginTx() (tx *sqlx.Tx, err error) {
	return o.DB.Beginx()
}

type OrderRepository interface {
	CreateOrder(tx *sqlx.Tx, ctx context.Context, form Order) (orderID int64, err error)
	CreateOrderDetails(tx *sqlx.Tx, ctx context.Context, form []OrderDetail) (err error)
	BeginTx() (tx *sqlx.Tx, err error)
}
