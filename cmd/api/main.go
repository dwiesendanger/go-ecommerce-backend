package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"ecommerce-platform/internal/handlers"
	"ecommerce-platform/internal/services"
	"ecommerce-platform/internal/workers"
	"ecommerce-platform/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// DB & Redis Connect
	db, err := database.ConnectPostgres()
	if err != nil {
		log.Fatal(err)
	}
	rdb := database.ConnectRedis()

	// Worker Setup
	var wg sync.WaitGroup
	orderJobChan := workers.StartOrderWorkers(100, 3, &wg)

	// Services
	authService := services.NewAuthService(db)
	checkoutService := services.NewCheckoutService(db, orderJobChan)
	cartService := services.NewCartService(db)
	productService := services.NewProductService(db, rdb)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	cartHandler := handlers.NewCartHandler(cartService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(checkoutService)

	// ROUTES

	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		// PUBLIC ROUTES
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)
		v1.GET("/products", productHandler.GetAllProducts)

		// PROTECTED ROUTES
		protected := v1.Group("/")
		protected.Use(handlers.AuthMiddleware())
		{
			protected.POST("/products", productHandler.CreateProduct)

			// CART
			protected.GET("/cart", cartHandler.GetCart)
			protected.POST("/cart/items", cartHandler.AddItem)

			// ORDERS
			protected.POST("/orders", orderHandler.CreateOrder) // Checkout
			protected.GET("/orders", orderHandler.GetOrders)    // History
		}
	}

	// Server Setup with Graceful Shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Closing worker channel...")
	close(orderJobChan) // Signal workers to stop

	log.Println("Waiting for workers to finish tasks...")
	wg.Wait() // Wait for all workers to finish

	log.Println("Server exiting")
}
