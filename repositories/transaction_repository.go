package repositories

import (
	"backend-hanssen-hilman/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	GetTransactionByID(id int64) (*models.TransactionResponse, error)
	ListTransactionsByMerchantID(merchantId int64, limit, page int) ([]models.TransactionResponse, int64, error)
	ListTransactionsByCustomerID(customerId int64, limit, page int) ([]models.TransactionResponse, int64, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) GetTransactionByID(id int64) (*models.TransactionResponse, error) {
	var transaction models.TransactionResponse
	err := r.db.Model(&models.Transaction{}).
		Select("transactions.id, transactions.product_id, products.name as product_name, transactions.quantity, transactions.total_price, users.name as customer").
		Joins("left join products on transactions.product_id = products.id").
		Joins("left join users on transactions.customer_id = users.id").
		First(&transaction, "transactions.id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) ListTransactionsByMerchantID(merchantId int64, limit, page int) ([]models.TransactionResponse, int64, error) {
	var transactions []models.TransactionResponse
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(&models.Transaction{}).
		Select("transactions.id, transactions.product_id, products.name as product_name, transactions.quantity, transactions.total_price, users.name as customer").
		Joins("left join products on transactions.product_id = products.id").
		Joins("left join (?) as users on transactions.customer_id = users.id", r.db.Model(&models.User{}).Where("role = 'customer'"))

	err := query.
		Limit(limit).Offset(offset).Where("products.merchant_id = ?", merchantId).
		Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}
	return transactions, total, nil
}

func (r *transactionRepository) ListTransactionsByCustomerID(customerId int64, limit, page int) ([]models.TransactionResponse, int64, error) {
	var transactions []models.TransactionResponse
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(&models.Transaction{}).
		Select("transactions.id, transactions.product_id, products.name as product_name, transactions.quantity, transactions.total_price, users.name as merchant, transactions.created_at, transactions.updated_at").
		Joins("left join products on transactions.product_id = products.id").
		Joins("left join (?) as users on products.merchant_id = users.id", r.db.Model(&models.User{}).Where("role = 'merchant'"))

	err := query.
		Limit(limit).Offset(offset).Where("transactions.customer_id = ?", customerId).
		Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}
	return transactions, total, nil
}
