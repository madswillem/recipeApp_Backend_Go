package database

import (
	"os" 

	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB() *gorm.DB {
	dsn := os.Getenv("DB")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil || db == nil {
		panic("Error ")
	}
	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.RecipeSchema{})
	db.AutoMigrate(&models.RecipeGroupSchema{})
	db.AutoMigrate(&models.Avrg{})
	db.AutoMigrate(&models.IngredientsSchema{})
	db.AutoMigrate(&models.RatingStruct{})
	db.AutoMigrate(&models.NutritionalValue{})
	db.AutoMigrate(&models.DietSchema{})
	db.AutoMigrate(&models.IngredientDBSchema{})
}
