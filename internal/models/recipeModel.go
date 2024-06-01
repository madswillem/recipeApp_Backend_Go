package models

import (
	"errors"
	"math"
	"net/http"

	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/tools"
	"gorm.io/gorm"
)

type RecipeSchema struct {
	BaseModel
	Title            string              `json:"title"`
	Ingredients      []IngredientsSchema `json:"ingredients" gorm:"foreignKey:RecipeSchemaID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Preparation      string              `json:"preparation"`
	Cuisine          string              `json:"cuisine"`
	CookingTime      int                 `json:"cookingtime"`
	Image            string              `json:"image"`
	NutriScore       string              `json:"nutriscore"`
	NutritionalValue NutritionalValue    `json:"nutritional_value" gorm:"polymorphic:Owner; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Diet             DietSchema          `json:"diet" gorm:"polymorphic:Owner; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Selected         int                 `json:"selected"`
	Rating           RatingStruct        `json:"rating" gorm:"polymorphic:Owner; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Version          int                 `json:"__v"`
	RecipeGroup	 []*RecipeGroupSchema`gorm:"many2many:recipe_recipegroups"` 
}

func (recipe *RecipeSchema) CheckIfExistsByTitle(db *gorm.DB) (bool, *error_handler.APIError) {
	var result struct {
		Found bool
	}
	err := db.Raw("SELECT EXISTS(SELECT * FROM recipe_schemas WHERE title = ?) AS found;", recipe.Title).Scan(&result).Error
	return result.Found, error_handler.New("database error", http.StatusInternalServerError, err)
}

func (recipe *RecipeSchema) CheckIfExistsByID(db *gorm.DB) (bool, *error_handler.APIError) {
	var result struct {
		Found bool
	}
	err := db.Raw("SELECT EXISTS(SELECT * FROM recipe_schemas WHERE id = ?) AS found;", recipe.ID).Scan(&result).Error
	if err != nil {
		return false, error_handler.New("database error", http.StatusInternalServerError, err)
	}
	return result.Found, nil
}

func (recipe *RecipeSchema) GetRecipeByID(db *gorm.DB ,reqData map[string]bool) *error_handler.APIError {
	req := db
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

func (recipe *RecipeSchema) AddNutritionalValue(db *gorm.DB) *error_handler.APIError {
	for _, ingredient := range recipe.Ingredients {
		var nutritionalValue NutritionalValue
		err := db.Joins("JOIN ingredients_schemas ON nutritional_values.owner_id = ingredients_schemas.id").
			Where("ingredients_schemas.ingredient = ?", ingredient.Ingredient).
			First(&nutritionalValue).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) && !ingredient.NutritionalValue.Edited {
				return error_handler.New("ingredient not found please add nutritional value and set edited to true", http.StatusBadRequest, err)
			} else if errors.Is(err, gorm.ErrRecordNotFound) && ingredient.NutritionalValue.Edited {
				err := ingredient.createIngredientDBEntry(db)
				if err != nil {
					return err
				}
			} else {
				return error_handler.New("database error", http.StatusInternalServerError, err)
			}
		} else {
			if ingredient.NutritionalValue.Edited {
				return error_handler.New("ingredient already exists", http.StatusBadRequest, err)
			} else if !ingredient.NutritionalValue.Edited {
				ingredient.NutritionalValue = nutritionalValue
			}
		}
	}
	return nil
}

func (recipe *RecipeSchema) UpdateSelected(change int, user *UserModel, db *gorm.DB) *error_handler.APIError {
	apiErr := recipe.GetRecipeByID(db ,map[string]bool{"everything": true})
	if apiErr != nil {
		return apiErr
	}
	recipe.Selected += change
	exists, apiErr := recipe.CheckIfExistsByID(db)
	if apiErr != nil {
		return apiErr
	}
	if !exists {
		return error_handler.New("recipe not found", http.StatusNotFound, gorm.ErrRecordNotFound)
	}

	apiErr = recipe.Rating.Update(change)
	if apiErr != nil {
		return apiErr
	}

	err := db.Save(recipe).Error
	if err != nil {
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}
	
	if user.ID == 0 {
		return nil
	}
	apiErr = user.AddRecipeToGroup(db, recipe)

	return apiErr
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

func (recipe *RecipeSchema) Create(db *gorm.DB) *error_handler.APIError {
	err := recipe.CheckForRequiredFields()
	if err != nil {
		return err
	}

	recipe.Rating.DefaultRatingStruct(recipe.Title)
	for i := 0; i < len(recipe.Ingredients); i++ {
		recipe.Ingredients[i].Rating.DefaultRatingStruct(recipe.Ingredients[i].Ingredient)
	}

	err = recipe.AddNutritionalValue(db)
	if err != nil {
		return err
	}

	return recipe.SubmitToDB(db)
}

func (recipe *RecipeSchema) GetSimilarityWithGroup(group RecipeGroupSchema) (float64, *error_handler.APIError) {
	//Ings
	sameIngs := make([]bool, len(recipe.Ingredients))
	sameAvrgIngs := make([]bool, len(group.AvrgIngredients))

	for i, ingredient := range recipe.Ingredients {
		for y, avrgIngredient := range group.AvrgIngredients {
			if ingredient.Ingredient == avrgIngredient.Name {
				sameIngs[i] = true
				sameAvrgIngs[y] = true
				break
			}
		}
	}

	var allIngs []float64
	for _, ing := range sameIngs {
		if !ing {
			allIngs = append(allIngs, 0)
		}
	}
	for i, avrging := range sameAvrgIngs {
		if avrging {
			allIngs = append(allIngs, group.AvrgIngredients[i].Percentige/float64(len(group.Recipes)))
		}
		if !avrging {
			allIngs = append(allIngs, 1-group.AvrgIngredients[i].Percentige/float64(len(group.Recipes)))
		}
	}
	simIngs := tools.CalculateAverage(allIngs)

	// Cuisine
	arrCuisine := make([]float64, len(group.AvrgCuisine))
	for i, cuisine := range group.AvrgCuisine {
		if cuisine.Name == recipe.Cuisine {
			arrCuisine[i] = cuisine.Percentige / float64(len(group.Recipes))
		}
		if cuisine.Name != recipe.Cuisine {
			arrCuisine[i] = 1 - cuisine.Percentige/float64(len(group.Recipes))
		}
	}
	simCuisine := tools.CalculateAverage(arrCuisine)
	// Diet
	sumDiet := 0.0
	switch {
	case !recipe.Diet.Vegetarien:
		sumDiet += group.AvrgVegetarien / float64(len(group.Recipes))
	case !recipe.Diet.Vegan:
		sumDiet += group.AvrgVegan / float64(len(group.Recipes))
	case !recipe.Diet.LowCal:
		sumDiet += group.AvrgLowCal / float64(len(group.Recipes))
	case !recipe.Diet.LowCarb:
		sumDiet += group.AvrgLowCarb / float64(len(group.Recipes))
	case !recipe.Diet.Keto:
		sumDiet += group.AvrgKeto / float64(len(group.Recipes))
	case !recipe.Diet.Paleo:
		sumDiet += group.AvrgPaleo / float64(len(group.Recipes))
	case !recipe.Diet.LowFat:
		sumDiet += group.AvrgLowFat / float64(len(group.Recipes))
	case !recipe.Diet.FoodCombining:
		sumDiet += group.AvrgFoodCombining / float64(len(group.Recipes))
	case !recipe.Diet.WholeFood:
		sumDiet += group.AvrgWholeFood / float64(len(group.Recipes))
	}
	sim := sumDiet
	sim += simIngs * 5.0
	sim += simCuisine * 4
	sim /= 10.0

	if math.IsNaN(sim) {
		return 0.0, error_handler.New("Similarity is not a number", http.StatusInternalServerError, nil)
	}

	return sim, nil
}

