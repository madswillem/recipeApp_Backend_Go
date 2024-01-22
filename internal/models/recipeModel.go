package models

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/initializers"
	"gorm.io/gorm"
)

type RecipeSchema struct {
	ID               uint                `json:"_id" gorm:"primarykey"`
	Title            string              `json:"title"`
	Ingredients      []IngredientsSchema `json:"ingredients" gorm:"foreignKey:RecipeSchemaID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Preparation      string              `json:"preparation"`
	CookingTime      int                 `json:"cookingtime"`
	Image            string              `json:"image"`
	NutriScore       string              `json:"nutriscore"`
	NutritionalValue NutritionalValue    `json:"nutritional_value" gorm:"polymorphic:Owner; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Diet             DietSchema          `json:"diet" gorm:"polymorphic:Owner; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Selected         int                 `json:"selected"`
	Rating           RatingStruct        `json:"rating" gorm:"polymorphic:Owner; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Version          int                 `json:"__v"`
}

func (recipe *RecipeSchema) Delete(c *gin.Context) *error_handler.APIError {
	exists, apiErr := recipe.CheckIfExistsByID()
	if apiErr != nil {
		return apiErr
	}
	if !exists {
		return error_handler.New("recipe not found", http.StatusNotFound, gorm.ErrRecordNotFound)
	}

	err := initializers.DB.Delete(&recipe).Error
	if err != nil {
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}
	return nil
}

func (recipe *RecipeSchema) Update() *error_handler.APIError {
	err := initializers.DB.Updates(&recipe).First(&recipe).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return error_handler.New("recipe not found", http.StatusNotFound, gorm.ErrRecordNotFound)
		} else {
			return error_handler.New("database error", http.StatusInternalServerError, err)
		}
	}
	return nil
}

func (recipe *RecipeSchema) CheckIfExistsByTitle() (bool, *error_handler.APIError) {
	var result struct {
		Found bool
	}
	err := initializers.DB.Raw("SELECT EXISTS(SELECT * FROM recipe_schemas WHERE title = ?) AS found;", recipe.Title).Scan(&result).Error
	return result.Found, error_handler.New("database error", http.StatusInternalServerError, err)
}

func (recipe *RecipeSchema) CheckIfExistsByID() (bool, *error_handler.APIError) {
	var result struct {
		Found bool
	}
	err := initializers.DB.Raw("SELECT EXISTS(SELECT * FROM recipe_schemas WHERE id = ?) AS found;", recipe.ID).Scan(&result).Error
	if err != nil {
		return false, error_handler.New("database error", http.StatusInternalServerError, err)
	}
	return result.Found, nil
}

func (recipe *RecipeSchema) GetRecipeByID(reqData map[string]bool) *error_handler.APIError {
	req := initializers.DB
	if reqData["ingredients"] || reqData["everything"] {
		req = req.Preload("Ingredients")
	}
	if reqData["ingredient_nutri"] || reqData["everything"] {
		req = req.Preload("Ingredients.NutritionalValue")
	}
	if reqData["ingredient_rate"] || reqData["everything"] {
		req = req.Preload("Ingredients.Rating")
	}
	if reqData["rating"] || reqData["everything"] {
		req = req.Preload("Rating")
	}
	if reqData["nutritionalvalue"] || reqData["everything"] {
		req = req.Preload("NutritionalValue")
	}
	if reqData["diet"] || reqData["everything"] {
		req = req.Preload("Diet")
	}
	err := req.First(&recipe, "ID = ?", recipe.ID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return error_handler.New("recipe not found", http.StatusNotFound, err)
		} else {
			return error_handler.New("database error", http.StatusInternalServerError, err)
		}
	}

	return nil
}

func (recipe *RecipeSchema) AddNutritionalValue() *error_handler.APIError {
	for _, ingredient := range recipe.Ingredients {
		var nutritionalValue NutritionalValue
		err := initializers.DB.Joins("JOIN ingredients_schemas ON nutritional_values.owner_id = ingredients_schemas.id").
			Where("ingredients_schemas.ingredient = ?", ingredient.Ingredient).
			First(&nutritionalValue).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) && !ingredient.NutritionalValue.Edited {
				return error_handler.New("ingredient not found please add nutritional value and set edited to true", http.StatusBadRequest, err)
			} else if errors.Is(err, gorm.ErrRecordNotFound) && ingredient.NutritionalValue.Edited {
				err := ingredient.createIngredientDBEntry()
				if err != nil {
					return err
				}
			} else {
				return error_handler.New("database error", http.StatusInternalServerError, err)
			}
		} else if err == nil {
			if ingredient.NutritionalValue.Edited {
				return error_handler.New("ingredient already exists", http.StatusBadRequest, err)
			} else if !ingredient.NutritionalValue.Edited {
				ingredient.NutritionalValue = nutritionalValue
			}
		}
	}
	return nil
}

func (recipe *RecipeSchema) UpdateSelected(change int, c *gin.Context) *error_handler.APIError {
	recipe.Selected += change
	apiErr := recipe.Rating.Update(change)
	if apiErr != nil {
		return apiErr
	}

	exists, apiErr := recipe.CheckIfExistsByID()
	if apiErr != nil {
		return apiErr
	}
	if !exists {
		return error_handler.New("recipe not found", http.StatusNotFound, gorm.ErrRecordNotFound)
	}

	err := initializers.DB.Save(recipe).Error
	if err != nil {
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}

	return nil
}

func (recipe *RecipeSchema) CheckForRequiredFields() *error_handler.APIError {
	if recipe.Title == "" {
		return error_handler.New("missing required field", http.StatusBadRequest, errors.New("missing recipe title"))
	}
	if recipe.Ingredients == nil {
		return error_handler.New("missing recipe ingredients", http.StatusBadRequest, errors.New("missing recipe ingredients"))
	}
	if recipe.Preparation == "" {
		return error_handler.New("missing recipe ingredients", http.StatusBadRequest, errors.New("missing recipe preparation"))
	}
	for _, ingredient := range recipe.Ingredients {
		err := ingredient.CheckForRequiredFields()
		if err != nil {
			return error_handler.New("missing required field", http.StatusBadRequest, err)
		}
	}

	return nil
}

func (recipe *RecipeSchema) Create() *error_handler.APIError {
	err := recipe.CheckForRequiredFields()
	if err != nil {
		return err
	}

	recipe.Rating.DefaultRatingStruct(recipe.Title)
	for i := 0; i < len(recipe.Ingredients); i++ {
		recipe.Ingredients[i].Rating.DefaultRatingStruct(recipe.Ingredients[i].Ingredient)
	}

	err = recipe.AddNutritionalValue()
	if err != nil {
		return err
	}

	tx := initializers.DB.Begin()

	if err := tx.Create(&recipe).Error; err != nil {
		tx.Rollback()
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}

	return err
}
