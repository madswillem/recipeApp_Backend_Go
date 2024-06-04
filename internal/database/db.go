package database

import (
	"errors"
	"net/http"
	"os"

	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB(conf *gorm.Config) *gorm.DB {
	
	dsn := os.Getenv("DB")
	db, err := gorm.Open(postgres.Open(dsn), conf)

	if err != nil || db == nil {
		panic("Error ")
	}
	return db
}

func SubmitToDB[T any](db *gorm.DB, entity *T) *error_handler.APIError {
	tx := db.Begin()

	if err := tx.Create(entity).Error; err != nil {
		tx.Rollback()
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}
	return nil
}

func Update[T any](db *gorm.DB, entity *T) *error_handler.APIError {
	err := db.Updates(entity).First(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_handler.New("entry not found", http.StatusNotFound, gorm.ErrRecordNotFound)
		} else {
			return error_handler.New("database error", http.StatusInternalServerError, err)
		}
	}
	return nil
}

func Delete[T any](db *gorm.DB, entity *T) *error_handler.APIError {
	err := db.Delete(entity).Error
	if err != nil {
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}
	return nil
}
