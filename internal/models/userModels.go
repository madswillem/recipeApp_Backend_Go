package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/madswillem/recipeApp_Backend_Go/internal/database"
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
func (user *UserModel) GetAllRecipeGroups(db *gorm.DB) ([]RecipeGroupSchema, *error_handler.APIError) {
	var groups []RecipeGroupSchema
	fmt.Println(user.ID)
	err := db.Select(&groups).Where("user_id = ?", user.ID).Scan(&groups).Error;
	if  err != nil {
	    return []RecipeGroupSchema{}, error_handler.New("Database error when fetching RecipeGroups i", http.StatusInternalServerError, err)
	}

	return groups, nil
}
func (user *UserModel) AddRecipeToGroup(db *gorm.DB ,recipe *RecipeSchema) *error_handler.APIError {	
	groups, err := user.GetAllRecipeGroups(db)
	if err != nil {
		return err
	}
	if len(groups) < 1 {
		user.RecipeGroups = append(user.RecipeGroups, GroupNew(recipe))
		return database.Update(db, user)
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
		sortedGroups[0].Group.AddRecipeToGroup(recipe, db)
		user.RecipeGroups = append(user.RecipeGroups, GroupNew(recipe))
		return database.Update(db, user)
	}

	sortedGroups[0].Group.AddRecipeToGroup(recipe, db)
	return database.Update(db, user)
}
// Using db to extend an existing db like a search to show recipes similar to your intrests
func (user *UserModel) GetRecomendation(db *gorm.DB) (*error_handler.APIError, []RecipeSchema) {
	var recipes []RecipeSchema
	//Get recipes 
	db =	db.Joins("JOIN ingredients_schemas ON recipe_schemas.id = ingredients_schemas.recipe_schema_id").
		Joins("JOIN diet_schemas ON diet_schemas.owner_id = recipe_schemas.id").
		Group("recipe_schemas.id").
		Preload(clause.Associations).
		Preload("Ingredients.Rating").
		Preload("Ingredients.NutritionalValue")

	switch {
		case user.Settings.Diet.Vegetarien:
			db = db.Where("diet_schemas.vegetarien = ?", true)
		case user.Settings.Diet.Vegan:
			db = db.Where("diet_schemas.vegan = ?", true)
		case user.Settings.Diet.LowCal:
			db = db.Where("diet_schemas.lowcal = ?", true)
		case user.Settings.Diet.LowCarb:
			db = db.Where("diet_schemas.lowcarb = ?", true)
		case user.Settings.Diet.Keto:
			db = db.Where("diet_schemas.keto = ?", true)
		case user.Settings.Diet.Paleo:
			db = db.Where("diet_schemas.paleo = ?", true)
		case user.Settings.Diet.LowFat:
			db = db.Where("diet_schemas.lowfat = ?", true)
		case user.Settings.Diet.FoodCombining:
			db = db.Where("diet_schemas.food_combining = ?", true)
		case user.Settings.Diet.WholeFood:
			db = db.Where("diet_schemas.whole_food = ?", true)
	}
	err := db.Find(&recipes).Error
	if err != nil {
		return error_handler.New("Database error when fetching recipes", http.StatusInternalServerError, err), nil
	}

	groups, apiErr := user.GetAllRecipeGroups(db)
	if apiErr != nil {
		return apiErr, nil
	}

	var similarity []SimiliarityGroupRecipe
	for _, recipe := range recipes {
		for _, group := range groups {
			sim, apiErr := recipe.GetSimilarityWithGroup(group)
			if apiErr != nil {
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
