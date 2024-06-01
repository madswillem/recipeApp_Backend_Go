package models

import (
	"errors"
	"net/http"
	"time"

	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
func (b *BaseModel) Update(db *gorm.DB) *error_handler.APIError{
	err := db.Updates(&b).First(&b).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_handler.New("recipe not entry", http.StatusNotFound, gorm.ErrRecordNotFound)
		} else {
			return error_handler.New("database error", http.StatusInternalServerError, err)
		}
	}
	return nil

}
func (b *BaseModel) Delete(db *gorm.DB) *error_handler.APIError {
	err := db.Delete(&b).Error
	if err != nil {
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}
	return nil

}
func (b *BaseModel) SubmitToDB(db *gorm.DB) *error_handler.APIError {
	tx := db.Begin()

	if err := tx.Create(&b).Error; err != nil {
		tx.Rollback()
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}
	return nil
}
