package models

import (
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/tools"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserModel struct {
	ID 				string 				`database:"id"`
	LastLogin		time.Time 			`database:"last_login"`
	Cookie			string 				`database:"cookie"`
	IP				string 				`database:"ip"`
	RecipeGroups	[]RecipeGroupSchema `gorm:"foreignKey:UserID;"`
	Settings 		UserSettings 		`gorm:"embedded;embeddedPrefix:setting_"`
}
type UserSettings struct {
	Allergies	[]*IngredientDB `gorm:"many2many:user_allergies"`
	Diet		DietSchema `gorm:"polymorphic:Owner"`
}

func (user *UserModel) GetByCookie(db *gorm.DB) *error_handler.APIError{
	err := db.Preload(clause.Associations).First(&user, "Cookie = ?", user.Cookie).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return error_handler.New("user not found", http.StatusNotFound, err)
		} else {
			return error_handler.New("database error", http.StatusInternalServerError, err)
		}
	}

	return nil

}
func (user *UserModel) CheckIfExistsByCookie(db *sqlx.DB) bool {
	found := false
	err := db.Get(found ,"SELECT EXISTS(SELECT * FROM user_models WHERE Cookie = $1) AS found;", user.Cookie)
	if err != nil {
		error_handler.New("database error", http.StatusInternalServerError, err)
	}
	return found
}
func (user *UserModel) Create(db *sqlx.DB ,ip string) *error_handler.APIError{
	user.LastLogin = time.Now()
	user.IP = ip
	for {
		user.Cookie = tools.RandomString(20)
		if !user.CheckIfExistsByCookie(db) {
			break
		}
	}

	_, err := db.NamedExec(`INSERT INTO user (cookie, ip) VALUES (:cookie, :ip)`, user)
	if err != nil {
		return error_handler.New("Error inserting user: "+err.Error(), http.StatusInternalServerError, err)
	}
	
	return nil
}
// Using db to extend an existing db like a search to show recipes similar to your intrests
func (user *UserModel) GetRecomendation(db *gorm.DB) (*error_handler.APIError, []RecipeSchema) {
	return nil, nil
}
