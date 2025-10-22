package controllers

import (
	"backend-hanssen-hilman/models"
	"backend-hanssen-hilman/repositories"
	"backend-hanssen-hilman/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	transactionRepo repositories.TransactionRepository
	productRepo     repositories.ProductRepository
}

func NewTransactionController(transactionRepo repositories.TransactionRepository, productRepo repositories.ProductRepository) *TransactionController {
	return &TransactionController{
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
	}
}

func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	var req models.TransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerId := ctx.GetInt64("user_id")

	product, err := c.productRepo.GetProductByID(req.ProductId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		return
	}

	if product.Product.Quantity < req.Quantity {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient product quantity"})
		return
	}

	productToUpdate := product.Product

	var totalPrice int64
	deliveryFee := 5000

	if product.Price < 15000 {
		totalPrice = req.Quantity*int64(product.Product.Price) + int64(deliveryFee)
	} else if product.Price > 50000 {
		totalPrice = req.Quantity * int64(product.Product.Price-(product.Product.Price*10/100))
	}

	newTransaction := models.Transaction{
		ProductId:  req.ProductId,
		Quantity:   req.Quantity,
		TotalPrice: float64(totalPrice),
		CustomerId: customerId,
	}

	if err := c.transactionRepo.CreateTransaction(&newTransaction); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	productToUpdate.Quantity -= req.Quantity

	if err := c.productRepo.UpdateProduct(&productToUpdate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product quantity"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully", "transaction": newTransaction})
}

func (c *TransactionController) GetTransactionByID(ctx *gin.Context) {
	id := ctx.Param("id")
	transactionId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}
	transaction, err := c.transactionRepo.GetTransactionByID(transactionId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transaction"})
		return
	}
	ctx.JSON(http.StatusOK, transaction)
}

func (c *TransactionController) ListTransactionsByMerchantID(ctx *gin.Context) {
	var req models.TransactionRequest
	merchantId := ctx.GetInt64("user_id")
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"invalid query parameters: ": err.Error()})
		return
	}

	req.Page, req.Limit = util.SetPaginationDefaults(req.Page, req.Limit)

	transactions, totalRecords, err := c.transactionRepo.ListTransactionsByMerchantID(merchantId, req.Limit, req.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	totalPages := util.CalculateTotalPages(totalRecords, req.Limit)

	ctx.JSON(http.StatusOK, gin.H{
		"transactions":  transactions,
		"total_records": totalRecords,
		"current_page":  req.Page,
		"page_size":     req.Limit,
		"total_pages":   totalPages,
	})
}

func (c *TransactionController) ListTransactionsByCustomerID(ctx *gin.Context) {
	var req models.TransactionRequest
	customerId := ctx.GetInt64("user_id")
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"invalid query parameters: ": err.Error()})
		return
	}

	req.Page, req.Limit = util.SetPaginationDefaults(req.Page, req.Limit)

	transactions, totalRecords, err := c.transactionRepo.ListTransactionsByCustomerID(customerId, req.Limit, req.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	totalPages := util.CalculateTotalPages(totalRecords, req.Limit)

	ctx.JSON(http.StatusOK, gin.H{
		"transactions":  transactions,
		"total_records": totalRecords,
		"current_page":  req.Page,
		"page_size":     req.Limit,
		"total_pages":   totalPages,
	})
}
