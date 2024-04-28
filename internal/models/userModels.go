package models

import (
	"net/http"
	"time"

	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/initializers"
	"github.com/madswillem/recipeApp_Backend_Go/internal/tools"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserModel struct {
	gorm.Model
	LastLogin	time.Time
	RecipeGroups	RecipeGroupSchema `gorm:"foreignKey:UserID;"`
	Cookie		string
	IP		string
}

func (user *UserModel) GetByCookie() *error_handler.APIError{
	err := initializers.DB.Preload(clause.Associations).First(&user, "Cookie = ?", user.Cookie).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return error_handler.New("user not found", http.StatusNotFound, err)
		} else {
			return error_handler.New("database error", http.StatusInternalServerError, err)
		}
	}

	return nil

}
func CheckIfExistsByCookie(cookie string) bool {
	var result struct {
		Found bool
		Error error_handler.APIError
	}
	err := initializers.DB.Raw("SELECT EXISTS(SELECT * FROM user_models WHERE Cookie = ?) AS found;", cookie).Scan(&result).Error
	if err != nil {
		error_handler.New("database error", http.StatusInternalServerError, err)
	}
	return result.Found
}
func (user *UserModel) Create(ip string) *error_handler.APIError{
	user.LastLogin = time.Now()
	user.IP = ip
	for {
		user.Cookie = tools.RandomString(20)
		if !CheckIfExistsByCookie(user.Cookie) {
			break
		}
	}

	tx := initializers.DB.Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}
	
	return nil
} 
