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
// Using query to extend an existing query like a search to show recipes similar to your intrests
func (user *UserModel) GetRecomendation(query *gorm.DB) (*error_handler.APIError, []RecipeSchema) {
	var recipes []RecipeSchema
	//Get recipes 
	query =	query.Joins("JOIN ingredients_schemas ON recipe_schemas.id = ingredients_schemas.recipe_schema_id").
		Joins("JOIN diet_schemas ON diet_schemas.owner_id = recipe_schemas.id").
		Group("recipe_schemas.id").
		Preload(clause.Associations).
		Preload("Ingredients.Rating").
		Preload("Ingredients.NutritionalValue")

	switch {
		case user.Settings.Diet.Vegetarien:
			query = query.Where("diet_schemas.vegetarien = ?", true)
		case user.Settings.Diet.Vegan:
			query = query.Where("diet_schemas.vegan = ?", true)
		case user.Settings.Diet.LowCal:
			query = query.Where("diet_schemas.lowcal = ?", true)
		case user.Settings.Diet.LowCarb:
			query = query.Where("diet_schemas.lowcarb = ?", true)
		case user.Settings.Diet.Keto:
			query = query.Where("diet_schemas.keto = ?", true)
		case user.Settings.Diet.Paleo:
			query = query.Where("diet_schemas.paleo = ?", true)
		case user.Settings.Diet.LowFat:
			query = query.Where("diet_schemas.lowfat = ?", true)
		case user.Settings.Diet.FoodCombining:
			query = query.Where("diet_schemas.food_combining = ?", true)
		case user.Settings.Diet.WholeFood:
			query = query.Where("diet_schemas.whole_food = ?", true)
	}
	err := query.Find(&recipes).Error
	if err != nil {
		return error_handler.New("Database error", http.StatusInternalServerError, err), nil
	}

	groups, apiErr := user.GetAllRecipeGroups()
	if apiErr != nil {
		return apiErr, nil
	}

	var similarity []SimiliarityGroupRecipe
	for _, recipe := range recipes {
		for _, group := range groups {
			sim, apiErr := recipe.GetSimilarityWithGroup(group)
			if err != nil {
				return apiErr, nil
			}
			similarity = append(similarity, SimiliarityGroupRecipe{Recipe: recipe, Group: group, Similarity: sim})
		}	
	}
	
	for i, sim := range similarity {
		if sim.Similarity < 80.0 {
			similarity[i] = similarity[len(similarity)-1]
			similarity = similarity[:len(similarity)-1]
		}
	}

	
	return nil, recipes
}
