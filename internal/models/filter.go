package models

import (
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"gorm.io/gorm"
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
