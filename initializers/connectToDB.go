package initializers

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"rezeptapp.ml/goApp/models"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil || DB == nil {
		panic("Error ")
	}

	DB.AutoMigrate(&models.RecipeSchema{})
	DB.AutoMigrate(&models.IngredientsSchema{})
	DB.AutoMigrate(&models.RatingStruct{})
}
