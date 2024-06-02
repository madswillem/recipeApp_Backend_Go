package database

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectToDB() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
		    SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
		    LogLevel:                  logger.Info,            // Log level
		    IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
		    Colorful:                  true,                   // Disable color
		},
	)
	dsn := os.Getenv("DB")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

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
