package models

import (
	"errors"
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
	RecipeGroups	[]RecipeGroupSchema `gorm:"foreignKey:UserID;"`
	Settings 	UserSettings `gorm:"embedded;embeddedPrefix:setting_"`
	Cookie		string
	IP		string
}
type UserSettings struct {
	Allergies	[]*IngredientDBSchema `gorm:"many2many:user_allergies"`
	Diet		DietSchema `gorm:"polymorphic:Owner"`
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
func (user *UserModel) Update() *error_handler.APIError{
	err := initializers.DB.Updates(&user).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_handler.New("recipe not found", http.StatusNotFound, gorm.ErrRecordNotFound)
		} else {
			return error_handler.New("database error", http.StatusInternalServerError, err)
		}
	}
	return nil

}
func (user *UserModel) GetAllRecipeGroups() ([]RecipeGroupSchema, *error_handler.APIError) {
	var group []RecipeGroupSchema
	if err := initializers.DB.Preload(clause.Associations).Find(&group, "user_id = ?", user.ID).Error; err != nil {
	    return []RecipeGroupSchema{}, error_handler.New("Database error", http.StatusInternalServerError, err)
	}

	return group, nil
}
func (user *UserModel) AddRecipeToGroup(recipe *RecipeSchema) *error_handler.APIError {	
	groups, err := user.GetAllRecipeGroups()
	if err != nil {
		return err
	}
	if len(groups) < 1 {
		user.RecipeGroups = append(user.RecipeGroups, GroupNew(recipe))
		return user.Update()
	}
	sortedGroups := make([]SimiliarityGroupRecipe, len(groups))

	for num, group := range groups {
		sortedGroups[num].Group = group
		sortedGroups[num].Similarity, err = recipe.GetSimilarityWithGroup(group)
		if err != nil {
			return err
		}
	}

	sortedGroups = SortSimilarity(sortedGroups)

	if sortedGroups[0].Similarity <= 90.0 {
		sortedGroups[0].Group.AddRecipeToGroup(recipe)
		user.RecipeGroups = append(user.RecipeGroups, GroupNew(recipe))
		return user.Update()
	}

	sortedGroups[0].Group.AddRecipeToGroup(recipe)
	return user.Update()
}
