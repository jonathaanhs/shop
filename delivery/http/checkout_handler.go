package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/learn/shop/repository"
	"github.com/learn/shop/usecase"
)

type CheckoutHandler struct {
	CheckoutUsecase usecase.CheckoutUsecase
}

func NewCheckoutHandler(checkoutUsecase usecase.CheckoutUsecase) CheckoutHandler {
	return CheckoutHandler{CheckoutUsecase: checkoutUsecase}
}

func (ch CheckoutHandler) InitRouter(router *gin.Engine) {
	router.POST("/checkout", ch.Checkout)
}

func (ch CheckoutHandler) Checkout(c *gin.Context) {
	var input []repository.OrderDetail

	c.BindJSON(&input)

	res, err := ch.CheckoutUsecase.Checkout(c.Request.Context(), input)
	if err != nil {
		log.Printf("[handler][Checkout] error while do Checkout %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Internal Server Error", "InternalMessage": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Success", "Response": res})
	c.Abort()
}
