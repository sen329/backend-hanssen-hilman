package models

import "time"

type Product struct {
	Id          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	Price       float64   `gorm:"column:price" json:"price"`
	MerchantId  int64     `gorm:"column:merchant_id" json:"merchant_id"`
	Quantity    int64     `gorm:"column:quantity" json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductDetail struct {
	Product
	MerchantName string `gorm:"column:merchant_name" json:"merchant_name"`
}

type ProductRequest struct {
	Name         string  `form:"name" json:"name"`
	Description  string  `form:"description" json:"description"`
	MerchantName string  `form:"merchant_name" json:"merchant_name"`
	Price        float64 `form:"price" json:"price"`
	MinPrice     float64 `form:"min_price" json:"min_price"`
	MaxPrice     float64 `form:"max_price" json:"max_price"`
	Quantity     int64   `form:"quantity" json:"quantity"`
	Page         int     `form:"page"`
	Limit        int     `form:"limit"`
}

type ProductResponse struct {
	Id           int64   `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	MerchantName string  `json:"merchant_name"`
	Quantity     int64   `json:"quantity"`
}

type PaginatedProductResponse struct {
	Products     []ProductResponse `json:"products"`
	TotalRecords int64             `json:"total_records"`
	CurrentPage  int               `json:"current_page"`
	PageSize     int               `json:"page_size"`
	TotalPages   int               `json:"total_pages"`
}
