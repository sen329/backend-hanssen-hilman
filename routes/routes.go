package routes

import (
	"backend-hanssen-hilman/controllers"
	"backend-hanssen-hilman/database"
	"backend-hanssen-hilman/repositories"
	"backend-hanssen-hilman/routes/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")

	userController := controllers.NewUserController(repositories.NewUserRepository(database.DB))

	// User Routes
	userRoutes := v1.Group("/users")
	{
		userRoutes.POST("/login", userController.Login)
		userRoutes.POST("/register", userController.Register)
	}

	// Product Routes
	productController := controllers.NewProductController(repositories.NewProductRepository(database.DB))
	productMerchantRoutes := v1.Group("/product/merchant")
	productMerchantRoutes.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("merchant"))
	{
		productMerchantRoutes.POST("/", productController.CreateProduct)
		productMerchantRoutes.PUT("/:id", productController.UpdateProduct)
		productMerchantRoutes.DELETE("/:id", productController.DeleteProduct)
		productMerchantRoutes.GET("/", productController.GetProductsByMerchantID)
		productMerchantRoutes.GET("/:id", productController.GetProductByID)
	}

	productRoutes := v1.Group("/products")
	productRoutes.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("customer"))
	{
		productRoutes.GET("/", productController.ListProducts)
		productRoutes.GET("/:id", productController.GetProductByID)
	}

	transactionController := controllers.NewTransactionController(repositories.NewTransactionRepository(database.DB), repositories.NewProductRepository(database.DB))

	// Merchant Routes
	merchantTransactionRoutes := v1.Group("/transactions/merchant")
	merchantTransactionRoutes.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("merchant"))
	{
		merchantTransactionRoutes.GET("/:id", transactionController.GetTransactionByID)
		merchantTransactionRoutes.GET("/", transactionController.ListTransactionsByMerchantID)
	}

	// Customer Routes
	customerTransactionRoutes := v1.Group("/transactions/customer")
	customerTransactionRoutes.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("customer"))
	{
		customerTransactionRoutes.POST("/", transactionController.CreateTransaction)
		customerTransactionRoutes.GET("/", transactionController.ListTransactionsByCustomerID)
		customerTransactionRoutes.GET("/:id", transactionController.GetTransactionByID)
	}

	// Run the server
	router.Run(":" + os.Getenv("PORT"))
	return router
}
