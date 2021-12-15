package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type promoRepo struct {
	DB *sqlx.DB
}

type Promo struct {
	PromoID   int64   `json:"promo_id" db:"promo_id"`
	ProductID int64   `json:"product_id" db:"product_id"`
	PromoType string  `json:"promo_type" db:"promo_type"`
	Reward    float64 `json:"reward" db:"reward"`
	MinQty    int64   `json:"min_qty" db:"min_qty"`
}

func NewPromoRepository(db *sqlx.DB) PromoRepository {
	return promoRepo{
		DB: db,
	}
}

func (o promoRepo) GetPromoByProductID(ctx context.Context, productID int64) (res Promo, err error) {
	rows, err := o.DB.QueryxContext(ctx, "select promo_id, product_id, promo_type, reward, min_qty from promos where product_id = $1", productID)
	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&res)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

func (o promoRepo) GetAllPromo(ctx context.Context) (res []Promo, err error) {
	rows, err := o.DB.QueryxContext(ctx, "select promo_id, product_id, promo_type, reward, min_qty from promos order by promo_id asc")
	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := Promo{}
		err = rows.StructScan(&tmp)
		if err != nil {
			return res, err
		}

		res = append(res, tmp)
	}

	return res, nil
}

type PromoRepository interface {
	GetPromoByProductID(ctx context.Context, productID int64) (res Promo, err error)
	GetAllPromo(ctx context.Context) (res []Promo, err error)
}
