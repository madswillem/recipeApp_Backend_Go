package models

import (
	"net/http"

	"github.com/madswillem/recipeApp_Backend_Go/internal/database"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"gorm.io/gorm"
)

type RecipeGroupSchema struct {
	gorm.Model
	UserID		  	  uint
	Recipes           []*RecipeSchema  `gorm:"many2many:recipe_recipegroups"` 
	AvrgIngredients   []Avrg `gorm:"foreignKey:GroupID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AvrgCuisine       []Avrg `gorm:"foreignKey:GroupID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AvrgVegetarien    float64
	AvrgVegan         float64
	AvrgLowCal        float64
	AvrgLowCarb       float64
	AvrgKeto          float64
	AvrgPaleo         float64
	AvrgLowFat        float64
	AvrgFoodCombining float64
	AvrgWholeFood     float64
}

type Avrg struct {
	ID              uint `gorm:"primarykey"`
	GroupID		uint
	Name            string
	Percentige      float64
}

func (group *RecipeGroupSchema) GetRecipeGroupByID(db *gorm.DB ,reqData map[string]bool) *error_handler.APIError {
	err := db.Find(&group, "ID = ?", group.ID).Error
	if err != nil {
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}
	return nil
}
func GetAllRecipeGroups(db *gorm.DB) (*error_handler.APIError, []RecipeGroupSchema) {
	var groups []RecipeGroupSchema
	err := db.Find(&groups).Error
	if err != nil {
		return error_handler.New("database error", http.StatusInternalServerError, err), nil
	}
	return nil, groups
}
func (group *RecipeGroupSchema) AddRecipeToGroup(recipe *RecipeSchema, db *gorm.DB) {
	for _, name := range recipe.Ingredients {
		added := false
		for _, avrgName := range group.AvrgIngredients {
			if name.Ingredient == avrgName.Name {
				avrgName.Percentige += 1
				added = true
			}
		}
		if !added {
			group.AvrgIngredients = append(group.AvrgIngredients, Avrg{Name: name.Ingredient, Percentige: 1})
		}
	}
	for _, cuisine := range group.AvrgCuisine {
		if recipe.Cuisine == cuisine.Name {
			cuisine.Percentige += 1
		}
	}
	switch {
	case recipe.Diet.Vegetarien:
		group.AvrgVegetarien += 1
	case recipe.Diet.Vegan:
		group.AvrgVegan += 1
	case recipe.Diet.LowCarb:
		group.AvrgLowCarb += 1
	case recipe.Diet.Keto:
		group.AvrgKeto += 1
	case recipe.Diet.Paleo:
		group.AvrgPaleo += 1
	case recipe.Diet.LowFat:
		group.AvrgLowFat += 1
	case recipe.Diet.FoodCombining:
		group.AvrgFoodCombining += 1
	case recipe.Diet.WholeFood:
		group.AvrgWholeFood += 1
	}
	group.Recipes = append(group.Recipes, recipe)
	database.Update(db, group)
}

func GroupNew(recipe *RecipeSchema) RecipeGroupSchema {
	new := RecipeGroupSchema{}
	new.Recipes = append(new.Recipes, recipe)
	for _, ing := range recipe.Ingredients {
		new.AvrgIngredients = append(new.AvrgIngredients, Avrg{Name: ing.Ingredient, Percentige: 1})
	}
	new.AvrgCuisine = append(new.AvrgCuisine, Avrg{Name: recipe.Cuisine, Percentige: 1})
	switch {
		case recipe.Diet.Vegetarien: new.AvrgVegetarien = 1
		case recipe.Diet.Vegan: new.AvrgVegan = 1
		case recipe.Diet.LowCal: new.AvrgLowCal = 1
		case recipe.Diet.LowCarb: new.AvrgLowCarb = 1
		case recipe.Diet.Keto: new.AvrgKeto = 1
		case recipe.Diet.Paleo: new.AvrgPaleo = 1
		case recipe.Diet.LowFat: new.AvrgLowFat = 1
		case recipe.Diet.FoodCombining: new.AvrgFoodCombining = 1
		case recipe.Diet.WholeFood: new.AvrgWholeFood = 1
	}

	return new
}
type SimiliarityGroupRecipe struct {
	Group RecipeGroupSchema
	Recipe RecipeSchema
	Similarity float64
}
func SortSimilarity( groups []SimiliarityGroupRecipe ) []SimiliarityGroupRecipe {
	len := len(groups)
	for i := 0; i < len-1; i++ {
		for j := 0; j < len-i-1; j++ {
			if groups[j].Similarity > groups[j+1].Similarity {
				groups[j], groups[j+1] = groups[j+1], groups[j]
			}
		}
	}
	return groups
}
