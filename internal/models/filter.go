package models

import (
	"net/http"

	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Filter struct {
	SearchText  string     `json:"search"`
	NutriScore  string     `json:"nutriscore"`
	CookingTime int        `json:"cookingtime"`
	Ingredients []string   `json:"ingredients"`
	Diet        DietSchema `json:"diet"`
}

func (f *Filter) AddFullTextSearchToQuery(query *gorm.DB) (*gorm.DB, *error_handler.APIError) {
	println(f.SearchText)
	query = query.Where("to_tsvector('english', recipe_schemas.title) @@ websearch_to_tsquery('english', ?)", f.SearchText).
		Or("to_tsvector('english', ingredients_schemas.ingredient) @@ websearch_to_tsquery('english', ?)", f.SearchText).
		Or("to_tsvector('english', recipe_schemas.preparation) @@ websearch_to_tsquery('english', ?)", f.SearchText).
		Order("COUNT(rating_structs.overall) desc")
	return query, nil
}

func (f *Filter) Filter(db *gorm.DB) (*[]RecipeSchema, *error_handler.APIError) {
	var recipes []RecipeSchema
	query := db.Joins("JOIN ingredients_schemas ON recipe_schemas.id = ingredients_schemas.recipe_schema_id").
		Joins("JOIN diet_schemas ON diet_schemas.owner_id = recipe_schemas.id").
		Group("recipe_schemas.id").
		Preload(clause.Associations).
		Preload("Ingredients.Rating").
		Preload("Ingredients.NutritionalValue")

	if f.Ingredients != nil {
		query = db.Where("ingredients_schemas.ingredient IN ?", f.Ingredients).
			Having("COUNT(DISTINCT ingredients_schemas.id) = ?", len(f.Ingredients))
	}

	switch {
	case f.Diet.Vegetarien:
		query = db.Where("diet_schemas.vegetarien = ?", true)
	case f.Diet.Vegan:
		query = db.Where("diet_schemas.vegan = ?", true)
	case f.Diet.LowCal:
		query = db.Where("diet_schemas.lowcal = ?", true)
	case f.Diet.LowCarb:
		query = db.Where("diet_schemas.lowcarb = ?", true)
	case f.Diet.Keto:
		query = db.Where("diet_schemas.keto = ?", true)
	case f.Diet.Paleo:
		query = db.Where("diet_schemas.paleo = ?", true)
	case f.Diet.LowFat:
		query = db.Where("diet_schemas.lowfat = ?", true)
	case f.Diet.FoodCombining:
		query = db.Where("diet_schemas.food_combining = ?", true)
	case f.Diet.WholeFood:
		query = db.Where("diet_schemas.whole_food = ?", true)
	case f.CookingTime > 0:
		query = db.Where("recipe_schemas.cooking_time <= ?", f.CookingTime)
	case f.NutriScore != "":
		query = db.Where("recipe_schemas.nutri_score = ?", f.NutriScore)
	}

	if f.SearchText != "" {
		query.Joins("JOIN rating_structs ON recipe_schemas.id = rating_structs.owner_id")
		query, _ = f.AddFullTextSearchToQuery(query)
	}

	err := query.Find(&recipes).Error

	if err != nil {
		return nil, error_handler.New("database error", http.StatusInternalServerError, err)
	}

	return &recipes, nil
}
