package models

import "time"

type Transaction struct {
	Id         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProductId  int64     `gorm:"column:product_id" json:"product_id"`
	Quantity   int64     `gorm:"column:quantity" json:"quantity"`
	TotalPrice float64   `gorm:"column:total_price" json:"total_price"`
	CustomerId int64     `gorm:"column:customer_id" json:"customer_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type TransactionRequest struct {
	ProductId  int64 `json:"product_id"`
	Quantity   int64 `json:"quantity"`
	CustomerId int64 `json:"customer_id"`
	Page       int   `form:"page"`
	Limit      int   `form:"limit"`
}

type TransactionResponse struct {
	Id          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProductId   int64     `gorm:"column:product_id" json:"product_id"`
	ProductName string    `gorm:"column:product_name" json:"product_name"`
	Quantity    int64     `gorm:"column:quantity" json:"quantity"`
	TotalPrice  float64   `gorm:"column:total_price" json:"total_price"`
	Customer    string    `gorm:"column:customer" json:"customer"`
	Merchant    string    `gorm:"column:merchant" json:"merchant"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PaginatedTransactionResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	TotalRecords int64                 `json:"total_records"`
	CurrentPage  int                   `json:"current_page"`
	PageSize     int                   `json:"page_size"`
	TotalPages   int                   `json:"total_pages"`
}
