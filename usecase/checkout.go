package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/learn/shop/repository"
)

type checkoutUsecase struct {
	OrderRepo   repository.OrderRepository
	ProductRepo repository.ProductRepository
	PromoRepo   repository.PromoRepository
}

type Checkout struct {
	Items       []string `json:"items"`
	TotalAmount float64  `json:"total_amount"`
}

func NewCheckoutUsecase(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, promoRepo repository.PromoRepository) CheckoutUsecase {
	return checkoutUsecase{
		OrderRepo:   orderRepo,
		ProductRepo: productRepo,
		PromoRepo:   promoRepo,
	}
}

func (c checkoutUsecase) Checkout(ctx context.Context, form []repository.OrderDetail) (res Checkout, err error) {
	tx, err := c.OrderRepo.BeginTx()
	if err != nil {
		log.Printf("error while do BeginTx %+v", err)
		return res, err
	}

	defer tx.Rollback()

	orderID, err := c.OrderRepo.CreateOrder(tx, ctx, repository.Order{
		Date:  time.Now(),
		Total: res.TotalAmount,
	})
	if err != nil {
		log.Printf("error while do CreateOrder %+v", err)
		return res, err
	}

	for i, v := range form {
		promo, err := c.PromoRepo.GetPromoByProductID(ctx, v.ProductID)
		if err != nil {
			log.Printf("error while do GetPromoByProductID %+v", err)
			return res, err
		}

		productDetail, err := c.ProductRepo.GetProductByProductID(ctx, v.ProductID)
		if err != nil {
			log.Printf("error while do GetProductByProductID %+v", err)
			return res, err
		}

		if productDetail.Qty < v.Qty {
			err = fmt.Errorf("the product %s qty is not enough to fulfill the request", productDetail.Name)
			return res, err
		}

		form[i].Price = float64(v.Qty) * productDetail.Price
		form[i].OrderID = orderID
		for i := 0; i < int(v.Qty); i++ {
			res.Items = append(res.Items, productDetail.Name)
		}

		if v.Qty >= promo.MinQty {
			switch promo.PromoType {
			case "product":
				if v.ProductID == int64(promo.Reward) {
					tmpQty := v.Qty - 1
					form[i].Price = float64(tmpQty) * productDetail.Price
				} else {
					productRewardDetail, err := c.ProductRepo.GetProductByProductID(ctx, int64(promo.Reward))
					if err != nil {
						log.Printf("error while do GetProductByProductID %+v", err)
						return res, err
					}
					res.Items = append(res.Items, productRewardDetail.Name)
				}
				form[i].PromoID = promo.PromoID
			case "discount":
				form[i].Price = (productDetail.Price * float64(v.Qty)) - ((productDetail.Price * float64(v.Qty)) * (promo.Reward / 100))
				form[i].PromoID = promo.PromoID
			}
		}

		err = c.ProductRepo.UpdateProductQtyByProductID(tx, ctx, repository.Product{
			ProductID: v.ProductID,
			Qty:       v.Qty,
		})
		if err != nil {
			log.Printf("error while do UpdateProductQtyByProductID %+v", err)
			return res, err
		}

		res.TotalAmount += form[i].Price
	}

	err = c.OrderRepo.CreateOrderDetails(tx, ctx, form)
	if err != nil {
		log.Printf("error while do CreateOrderDetails %+v", err)
		return res, err
	}

	tx.Commit()

	return res, nil
}

type CheckoutUsecase interface {
	Checkout(ctx context.Context, form []repository.OrderDetail) (res Checkout, err error)
}
