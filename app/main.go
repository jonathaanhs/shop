package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	checkoutHttpDelivery "github.com/learn/shop/delivery/http"
	"github.com/learn/shop/repository"
	"github.com/learn/shop/usecase"
	_ "github.com/lib/pq" //import for postgres driver
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")
	dbUser := os.Getenv("db_user")
	dbPass := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dsn := "postgresql://" + dbUser + ":" + dbPass + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 10)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	orderRepo := repository.NewOrderRepository(db)
	productRepo := repository.NewProductRepository(db)
	promoRepo := repository.NewPromoRepository(db)
	checkoutusecase := usecase.NewCheckoutUsecase(orderRepo, productRepo, promoRepo)

	checkoutHandler := checkoutHttpDelivery.NewCheckoutHandler(checkoutusecase)

	checkoutHandler.InitRouter(router)

	log.Println("Starting HTTP server")

	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
