package migrations

import (
	"backend-hanssen-hilman/models"
	"fmt"

	"gorm.io/gorm"
)

// Migrate runs the database migrations.
func Migrate(db *gorm.DB) {
	fmt.Println("Running migrations...")
	err := db.AutoMigrate(&models.User{}, &models.Product{}, &models.Transaction{})
	if err != nil {
		panic("failed to migrate database")
	}
	fmt.Println("Migrations completed successfully.")
}
