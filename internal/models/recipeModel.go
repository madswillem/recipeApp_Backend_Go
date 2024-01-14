package models

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/initializers"
)

type RecipeSchema struct {
	ID        		 uint 			 	 `json:"_id" gorm:"primarykey"`
	Title			 string				 `json:"title"` 
	Ingredients		 []IngredientsSchema `json:"ingredients" gorm:"foreignKey:RecipeSchemaID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Preparation		 string				 `json:"preparation"`
	CookingTime		 int				 `json:"cookingtime"`
	Image			 string				 `json:"image"`
	NutriScore		 string				 `json:"nutriscore"`
	NutritionalValue NutritionalValue    `json:"nutritional_value" gorm:"polymorphic:Owner; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Diet			 DietSchema			 `json:"diet" gorm:"polymorphic:Owner; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Selected		 int				 `json:"selected"`
	Rating			 RatingStruct		 `json:"rating" gorm:"polymorphic:Owner; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`	
	Version     	 int				 `json:"__v"`
}

func (recipe *RecipeSchema) Delete(c *gin.Context) error {
	exists, err := recipe.CheckIfExistsByID()
	if err != nil {
		return err
	}
	if !exists {
		return gorm.ErrRecordNotFound
	}

	err = initializers.DB.Delete(&recipe).Error
	return err
}

func (recipe *RecipeSchema) Update(c *gin.Context) error {
	err := initializers.DB.Updates(&recipe).Error
	return err
}

func (recipe *RecipeSchema) CheckIfExistsByTitle() (bool, error) {
	var result struct {
		Found bool
	}
	err := initializers.DB.Raw("SELECT EXISTS(SELECT * FROM recipe_schemas WHERE title = ?) AS found;", recipe.Title).Scan(&result).Error
	return result.Found, err
}

func (recipe *RecipeSchema) CheckIfExistsByID() (bool, error) {
	var result struct {
		Found bool
	}
	err := initializers.DB.Raw("SELECT EXISTS(SELECT * FROM recipe_schemas WHERE id = ?) AS found;", recipe.ID).Scan(&result).Error
	return result.Found, err
}

func (recipe *RecipeSchema) GetRecipeByID(c *gin.Context, reqData map[string]bool) error {
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
			error_handler.HandleError(c, http.StatusNotFound, "Recipe not found", err)
		} else {
			error_handler.HandleError(c, http.StatusInternalServerError, "Database error", err)
		}
	}

	return err
}

func (recipe *RecipeSchema) AddNutritionalValue() error {
	for _, ingredient := range recipe.Ingredients {
        var nutritionalValue NutritionalValue		
        err := initializers.DB.Joins("JOIN ingredients_schemas ON nutritional_values.owner_id = ingredients_schemas.id").
            Where("ingredients_schemas.ingredient = ?", ingredient.Ingredient).
            First(&nutritionalValue).Error
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) && !ingredient.NutritionalValue.Edited {
				return err
            } else if errors.Is(err, gorm.ErrRecordNotFound) && ingredient.NutritionalValue.Edited {
				err = ingredient.createIngredientDBEntry()
				if err != nil {
					return err
				}
            } else {
				return err
            }
        } else if err == nil {
            if ingredient.NutritionalValue.Edited {
				err = errors.New("ingredient already exists")
				return err
            } else if !ingredient.NutritionalValue.Edited {
                ingredient.NutritionalValue = nutritionalValue
            }
        }
    }
	return nil
}

func (recipe *RecipeSchema) UpdateSelected(change int, c *gin.Context) error {
	recipe.Selected += change
	recipe.Rating.Update(change, c)

	exists, err := recipe.CheckIfExistsByID()
	if err != nil {
		return err
	}
	if !exists {
		return gorm.ErrRecordNotFound
	}

	res := initializers.DB.Save(recipe)
	err = res.Error
	if err != nil {
		return err
	}

	return err
}

func (recipe *RecipeSchema) CheckForRequiredFields() error {
	if recipe.Title == "" {
		return errors.New("missing recipe Title")
	}
	if recipe.Ingredients == nil {
		return errors.New("missing recipe ingredients")
	}
	if recipe.Preparation == "" {
		return errors.New("missing recipe preparation")
	}
	for _, ingredient := range recipe.Ingredients {
		err := ingredient.CheckForRequiredFields()
		if err != nil {
			return err
		}
	}

	return nil
}

func (recipe *RecipeSchema) Create() error {
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
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return err
}