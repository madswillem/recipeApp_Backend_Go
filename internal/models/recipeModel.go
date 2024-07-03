package models

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/tools"
	"gorm.io/gorm"
)

type RecipeSchema struct {
	ID          string        `db:"id"`
	CreatedAt   time.Time     `db:"created_at"`
	Author      string        `db:"author"`
	Name        string        `db:"name"`
	Cuisine     string        `db:"cuisine"`
	Yield       int           `db:"yield"`
	YieldUnit   string        `db:"yield_unit"`
	PrepTime    string        `db:"prep_time"`
	CookingTime string        `db:"cooking_time"`
	Selected	int			  `db:"selected"`
	Version     int64         `db:"version"`
	Ingredients []IngredientsSchema 
	Diet 		DietSchema
	NutritionalValue	NutritionalValue
	Rating 				RatingStruct
	Steps 				[]StepsStruct
}

func (recipe *RecipeSchema) CheckIfExistsByTitle(db *gorm.DB) (bool, *error_handler.APIError) {
	var result struct {
		Found bool
	}
	err := db.Raw("SELECT EXISTS(SELECT * FROM recipe_schemas WHERE title = ?) AS found;", recipe.Name).Scan(&result).Error
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
	err := req.First(&recipe, "id = ?", recipe.ID).Error

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
			Where("ingredients_schemas.ingredient = ?", ingredient.Name).
			First(&nutritionalValue).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) && !ingredient.NutritionalValue.Edited {
				return error_handler.New(fmt.Sprintf("ingredient %s not found please add nutritional value and set edited to true", ingredient.Name), http.StatusBadRequest, err)
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
				return error_handler.New(fmt.Sprintf("Ingredient: %s already exists", ingredient.Name), http.StatusBadRequest, err)
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

	apiErr = recipe.Rating.Update(change)
	if apiErr != nil {
		return apiErr
	}
	fmt.Println(recipe.Rating.Overall)

	err := db.Model(&recipe).Update("selected", recipe.Selected).Error
	if err != nil {
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}
	err = db.Model(&recipe.Rating).Updates(recipe.Rating).Error
	if err != nil {
		return error_handler.New("database error", http.StatusInternalServerError, err)
	}
	
	//if user == nil {
	//	return nil
	//}
	//apiErr = user.AddRecipeToGroup(db, recipe)

	return nil
}

func (recipe *RecipeSchema) CheckForRequiredFields() *error_handler.APIError {
	if recipe.Name == "" {
		return error_handler.New("missing recipe name", http.StatusBadRequest, errors.New("missing recipe name"))
	}
	if recipe.Ingredients == nil {
		return error_handler.New("missing recipe ingredients", http.StatusBadRequest, errors.New("missing recipe ingredients"))
	}
	if recipe.Steps == nil {
		return error_handler.New("missing recipe steps", http.StatusBadRequest, errors.New("missing recipe steps"))
	}
	for _, ingredient := range recipe.Ingredients {
		err := ingredient.CheckForRequiredFields()
		if err != nil {
			return error_handler.New(fmt.Sprintf("missing required field in ingredient %s %s", ingredient.Name,err.Error()), http.StatusBadRequest, err)
		}
	}

	return nil
}

func (recipe *RecipeSchema) Create(db *sqlx.DB) *error_handler.APIError {
	apiErr := recipe.CheckForRequiredFields()
	if apiErr != nil {
		return apiErr
	}

	for i := 0; i < len(recipe.Ingredients); i++ {
		recipe.Ingredients[i].Rating.DefaultRatingStruct(nil, &recipe.Ingredients[i].ID)
	}

	//err = recipe.AddNutritionalValue(db)
	//if err != nil {
	//	return err
	//}

	tx := db.MustBegin()
	// Insert recipe
    query := `INSERT INTO recipes (author, name, cuisine, yield, yield_unit, prep_time, cooking_time, selected, version)
              VALUES (:author, :name, :cuisine, :yield, :yield_unit, :prep_time, :cooking_time, :selected, :version) RETURNING id`
    stmt, err := tx.PrepareNamed(query)
    if err != nil {
        return error_handler.New("Query error: " + err.Error(), http.StatusInternalServerError, err)
    }
    err = stmt.Get(&recipe.ID, recipe)
	stmt.Close()
    if err != nil {
		tx.Rollback()
        return error_handler.New("Dtabase error: " + err.Error(), http.StatusInternalServerError, err)
    }

	// Insert Rating
	recipe.Rating.DefaultRatingStruct(&recipe.ID, nil)
	query = `INSERT INTO rating (
				recipe_id, overall, mon, tue, wed, thu, fri, sat, sun, win, spr, sum, aut,
				thirtydegree, twentiedegree, tendegree, zerodegree, subzerodegree)
			VALUES (
				:recipe_id, :overall, :mon, :tue, :wed, :thu, :fri, :sat, :sun, :win, :spr, :sum, :aut,
				:thirtydegree, :twentiedegree, :tendegree, :zerodegree, :subzerodegree)`

	_, err = tx.NamedExec(query, recipe.Rating)
	if err != nil {
		tx.Rollback()
		return error_handler.New("Error inserting recipe: " + err.Error(), http.StatusInternalServerError, err)
	}

	// Insert Ingredient
	for _, ing := range recipe.Ingredients   {
		ing.RecipeID = recipe.ID
		err := ing.Create(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return error_handler.New("Error creating recipe", http.StatusInternalServerError, err)
	}


	fmt.Println(recipe.ID)

	return nil
}

func (recipe *RecipeSchema) GetSimilarityWithGroup(group RecipeGroupSchema) (float64, *error_handler.APIError) {
	//Ings
	sameIngs := make([]bool, len(recipe.Ingredients))
	sameAvrgIngs := make([]bool, len(group.AvrgIngredients))

	for i, ingredient := range recipe.Ingredients {
		for y, avrgIngredient := range group.AvrgIngredients {
			if ingredient.Name == avrgIngredient.Name {
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
