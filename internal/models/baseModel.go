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
	query	  *gorm.DB `gorm:"ignore"`
}
func (b BaseModel) Update() *error_handler.APIError{
	err := b.query.Updates(&b).First(&b).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_handler.New("recipe not entry", http.StatusNotFound, gorm.ErrRecordNotFound)
		} else {
			return error_handler.New("database error", http.StatusInternalServerError, err)
		}
	}
	return nil

}
func (b BaseModel) Delete() *error_handler.APIError {
	err := b.query.Delete(&b).Error
	if err != nil {
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}
	return nil

}