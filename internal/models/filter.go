package models

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
)

type Filter struct {
	SearchText  *string `db:"searchtext" json:"search"`
	NutriScore  *string `db:"nutriscore" json:"nutriscore"`
	Name        *string `db:"name" json:"name"`
	Cuisine     *string `db:"cuisine" json:"cuisine"`
	PrepTime    *string `db:"prep_time" json:"prep_time"`
	CookingTime *string `db:"cooking_time" json:"cooking_time"`
	Ingredients *[]string `json:"ingredients"`
	Diet        *DietSchema
}

func (f *Filter) Filter(db *sqlx.DB) (*[]RecipeSchema, *error_handler.APIError) {
	recipes := []RecipeSchema{}
	var where []string
	var args []interface{}

	if f.SearchText != nil {
		where = append(where, `to_tsvector('english', recipes.name) @@ websearch_to_tsquery('english', $1) 
					OR to_tsvector('english', ingredient.name) @@ websearch_to_tsquery('english', $1)
					OR to_tsvector('english', step.step) @@ websearch_to_tsquery('english', $1)`)
		args = append(args, f.SearchText)
	}
	if f.NutriScore != nil {
		args = append(args, f.NutriScore)
		where = append(where, fmt.Sprintf(`nutritional_value.nutriscore = :$%d`, len(args)))
	}
	if f.Cuisine != nil {
		args = append(args, f.Cuisine)
		where = append(where, fmt.Sprintf(`recipes.cuisine = $%d`, len(args)))
	}
	if f.PrepTime != nil {
		args = append(args, f.PrepTime)
		where = append(where, fmt.Sprintf(`recipes.prep_time <= $%d`, len(args)))
	}
	if f.CookingTime != nil {
		args = append(args, f.CookingTime)
		where = append(where, fmt.Sprintf(`recipes.cooking_time <= $%d`, len(args)))
	}
	if f.Ingredients != nil && len(*f.Ingredients) > 0 {
		for _, ing := range *f.Ingredients {
			args = append(args, ing)
			where = append(where, fmt.Sprintf(`ingredient.name = $%d`, len(args)))
		}
	}
	if f.Diet != nil {
		where = append(where, `AND diet.vegetarien = :vegetarien
								AND diet.vegan = :vegan
								AND diet.lowcal = :lowcal
								AND diet.lowcarb = :lowcarb
								AND diet.keto = :keto
								AND diet.paleo = :paleo
								AND diet.lowfat = :lowfat
								AND diet.food_combining = :food_combining
								AND diet.whole_food = :whole_food;`)
	}

	query := `SELECT DISTINCT recipes.*
				FROM recipes 
				LEFT JOIN recipe_ingredient ON recipes.id = recipe_ingredient.recipe_id
				LEFT JOIN ingredient ON ingredient.id = recipe_ingredient.ingredient_id 
				LEFT JOIN nutritional_value ON recipes.id = nutritional_value.recipe_id
				LEFT JOIN step ON recipes.id = step.recipe_id
				WHERE ` + strings.Join(where, " AND ")

	fmt.Println(query)

	err := db.Select(&recipes, query,args...)
	if err != nil {
		return nil, error_handler.New("Dtabase error: "+err.Error(), http.StatusInternalServerError, err)
	}

	return &recipes, nil
}
