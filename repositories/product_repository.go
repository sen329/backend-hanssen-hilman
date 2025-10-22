package repositories

import (
	"backend-hanssen-hilman/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id int64) (*models.ProductDetail, error)
	GetProductByMerchantID(id int64, page, limit int) ([]models.ProductDetail, int64, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id uint) error
	ListProducts(filter models.ProductRequest) ([]models.ProductDetail, int64, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetProductByID(id int64) (*models.ProductDetail, error) {
	var product models.ProductDetail
	err := r.db.Model(&models.Product{}).
		Select("products.*, users.name as merchant_name").
		Joins("left join users on products.merchant_id = users.id").
		First(&product, "products.id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetProductByMerchantID(id int64, page, limit int) ([]models.ProductDetail, int64, error) {
	var total int64
	var products []models.ProductDetail
	offset := (page - 1) * limit

	query := r.db.Model(&models.Product{}).
		Select("products.*, users.name as merchant_name").
		Joins("left join (?) as users on products.merchant_id = users.id", r.db.Model(&models.User{}).Where("role = 'merchant'"))

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).
		Find(&products, "products.merchant_id = ?", id).Error
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *productRepository) UpdateProduct(product *models.Product) error {
	return r.db.Model(product).Updates(product).Error
}

func (r *productRepository) DeleteProduct(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *productRepository) ListProducts(filter models.ProductRequest) ([]models.ProductDetail, int64, error) {
	var total int64
	var products []models.ProductDetail
	query := r.db.Model(&models.Product{}).
		Select("products.*, users.name as merchant_name").Debug().
		Joins("left join (?) as users on products.merchant_id = users.id", r.db.Model(&models.User{}).Where("role = 'merchant'"))

	if filter.Name != "" {
		query = query.Where("products.name LIKE ?", "%"+filter.Name+"%")
	}

	if filter.Description != "" {
		query = query.Where("products.description LIKE ?", "%"+filter.Description+"%")
	}

	if filter.MinPrice > 0 {
		query = query.Where("products.price >= ?", filter.MinPrice)
	}

	if filter.MaxPrice > 0 {
		query = query.Where("products.price <= ?", filter.MaxPrice)
	}

	if filter.Price > 0 {
		query = query.Where("products.price = ?", filter.Price)
	}

	if filter.MerchantName != "" {
		query = query.Where("users.name LIKE ?", "%"+filter.MerchantName+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.Limit(filter.Limit).Offset(offset).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil

}
