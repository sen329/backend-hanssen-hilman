package controllers

import (
	"backend-hanssen-hilman/models"
	"backend-hanssen-hilman/repositories"
	"backend-hanssen-hilman/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct {
	productRepo repositories.ProductRepository
}

func NewProductController(productRepo repositories.ProductRepository) *ProductController {
	return &ProductController{productRepo: productRepo}
}

func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req models.ProductRequest
	merchantId := ctx.GetInt64("user_id")
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newProduct := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		MerchantId:  merchantId,
		Quantity:    req.Quantity,
	}

	if err := c.productRepo.CreateProduct(&newProduct); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "product": newProduct})

}

func (c *ProductController) GetProductByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := c.productRepo.GetProductByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		}
		return
	}

	ctx.JSON(http.StatusOK, product)

}

func (c *ProductController) GetProductsByMerchantID(ctx *gin.Context) {
	var req models.ProductRequest
	merchantId := ctx.GetInt64("user_id")
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"invalid query parameters: ": err.Error()})
		return
	}

	req.Page, req.Limit = util.SetPaginationDefaults(req.Page, req.Limit)

	products, totalRecords, err := c.productRepo.GetProductByMerchantID(merchantId, req.Page, req.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list products"})
		return
	}

	totalPages := util.CalculateTotalPages(totalRecords, req.Limit)

	productResponses := []models.ProductResponse{}

	for _, p := range products {
		productRes := models.ProductResponse{
			Id:           p.Product.Id,
			Name:         p.Product.Name,
			Description:  p.Product.Description,
			Price:        p.Product.Price,
			MerchantName: p.MerchantName,
			Quantity:     p.Product.Quantity,
		}

		productResponses = append(productResponses, productRes)
	}

	ctx.JSON(http.StatusOK, models.PaginatedProductResponse{
		Products:     productResponses,
		TotalRecords: totalRecords,
		CurrentPage:  req.Page,
		PageSize:     req.Limit,
		TotalPages:   totalPages,
	})
}

func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	var req models.ProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := c.productRepo.GetProductByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		}
		return
	}

	productToUpdate := product.Product

	if req.Name != "" {
		productToUpdate.Name = req.Name
	}
	if req.Description != "" {
		productToUpdate.Description = req.Description
	}
	if req.Price > 0 {
		productToUpdate.Price = req.Price
	}
	if req.Quantity >= 0 {
		productToUpdate.Quantity = req.Quantity
	}

	if err := c.productRepo.UpdateProduct(&productToUpdate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := c.productRepo.DeleteProduct(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (c *ProductController) ListProducts(ctx *gin.Context) {
	var req models.ProductRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"invalid query parameters: ": err.Error()})
		return
	}

	req.Page, req.Limit = util.SetPaginationDefaults(req.Page, req.Limit)

	products, totalRecords, err := c.productRepo.ListProducts(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list products"})
		return
	}

	totalPages := util.CalculateTotalPages(totalRecords, req.Limit)

	productResponses := []models.ProductResponse{}

	for _, p := range products {
		productRes := models.ProductResponse{
			Id:           p.Id,
			Name:         p.Name,
			Description:  p.Description,
			Price:        p.Price,
			MerchantName: p.MerchantName,
			Quantity:     p.Quantity,
		}

		productResponses = append(productResponses, productRes)
	}

	ctx.JSON(http.StatusOK, models.PaginatedProductResponse{
		Products:     productResponses,
		TotalRecords: totalRecords,
		CurrentPage:  req.Page,
		PageSize:     req.Limit,
		TotalPages:   totalPages,
	})
}
