package models

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"gorm.io/gorm"
)

type RecipeSchema struct {
	ID               string    `db:"id"`
	CreatedAt        time.Time `db:"created_at"`
	Author           string    `db:"author"`
	Name             string    `db:"name"`
	Cuisine          string    `db:"cuisine"`
	Yield            int       `db:"yield"`
	YieldUnit        string    `db:"yield_unit"`
	PrepTime         string    `db:"prep_time"`
	CookingTime      string    `db:"cooking_time"`
	Selected         int       `db:"selected"`
	Version          int64     `db:"version"`
	Ingredients      []IngredientsSchema
	Diet             DietSchema
	NutritionalValue NutritionalValue
	Rating           RatingStruct `db:"rating"`
	Steps            []StepsStruct
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

func (recipe *RecipeSchema) GetRecipeByIDGORM(db *gorm.DB, reqData map[string]bool) *error_handler.APIError {
	return nil
}

func (recipe *RecipeSchema) GetRecipeByID(db *sqlx.DB) *error_handler.APIError {
	err := db.Get(recipe, `SELECT recipes.*,
								rt.id AS "rating.id", rt.created_at AS "rating.created_at",
								rt.recipe_id AS "rating.recipe_id", rt.overall AS "rating.overall", rt.mon AS "rating.mon",
								rt.tue AS "rating.tue", rt.wed AS "rating.wed", rt.thu AS "rating.thu", rt.fri AS "rating.fri",
								rt.sat AS "rating.sat", rt.sun AS "rating.sun", rt.win AS "rating.win",
								rt.spr AS "rating.spr", rt.sum AS "rating.sum", rt.aut AS "rating.aut",
								rt.thirtydegree AS "rating.thirtydegree", rt.twentiedegree AS "rating.twentiedegree",
								rt.tendegree AS "rating.tendegree", rt.zerodegree AS "rating.zerodegree",
								rt.subzerodegree AS "rating.subzerodegree"
							FROM recipes
							LEFT JOIN rating rt ON rt.recipe_id = recipes.id
							WHERE recipes.id = $1`, recipe.ID)
	if err != nil {
		return error_handler.New("An error ocurred fetching the recipe: "+err.Error(), http.StatusInternalServerError, err)
	}

	err = db.Select(&recipe.Steps, `SELECT * FROM step WHERE recipe_id = $1`, recipe.ID)
	if err != nil {
		return error_handler.New("An error ocurred fetching the steps: "+err.Error(), http.StatusInternalServerError, err)
	}

	err = db.Select(&recipe.Ingredients, `SELECT recipe_ingredient.*, ingredient.name AS name
										FROM recipe_ingredient
										INNER JOIN ingredient ON ingredient.id = recipe_ingredient.ingredient_id
										WHERE recipe_id = $1`, recipe.ID)
	if err != nil {
		return error_handler.New("An error ocurred fetching the ingredients: "+err.Error(), http.StatusInternalServerError, err)
	}

	return nil
}

func (recipe *RecipeSchema) UpdateSelected(change int, user *UserModel, db *sqlx.DB) *error_handler.APIError {
	apiErr := recipe.GetRecipeByID(db)
	if apiErr != nil {
		return apiErr
	}

	apiErr = recipe.Rating.Update(change)
	if apiErr != nil {
		return apiErr
	}

	fmt.Println(recipe.Rating.Overall)

	recipe.Selected += change
	recipe.Version += 1

	tx := db.MustBegin()
	tx.MustExec(`UPDATE "user" SET selected=$1, version=$2 WHERE id=$3`, recipe.Selected, recipe.Version, recipe.ID)
	query := `
    UPDATE weather_data
    SET
        overall = :overall,
        mon = :mon,
        tue = :tue,
        wed = :wed,
        thu = :thu,
        fri = :fri,
        sat = :sat,
        sun = :sun,
        win = :win,
        spr = :spr,
        sum = :sum,
        aut = :aut,
        thirtydegree = :thirtydegree,
        twentiedegree = :twentiedegree,
        tendegree = :tendegree,
        zerodegree = :zerodegree,
        subzerodegree = :subzerodegree
    WHERE some_condition = :some_condition`
	tx.NamedExec(query, recipe.Rating)
	err := tx.Commit()
	if err != nil {
		return error_handler.New("Error creating recipe", http.StatusInternalServerError, err)
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
			return error_handler.New(fmt.Sprintf("missing required field in ingredient %s %s", ingredient.Name, err.Error()), http.StatusBadRequest, err)
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
		return error_handler.New("Query error: "+err.Error(), http.StatusInternalServerError, err)
	}
	err = stmt.Get(&recipe.ID, recipe)
	stmt.Close()
	if err != nil {
		tx.Rollback()
		return error_handler.New("Dtabase error: "+err.Error(), http.StatusInternalServerError, err)
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
		return error_handler.New("Error inserting recipe: "+err.Error(), http.StatusInternalServerError, err)
	}

	// Insert Ingredient
	for _, ing := range recipe.Ingredients {
		ing.RecipeID = recipe.ID
		err := ing.Create(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	//Insert Steps
	for _, s := range recipe.Steps {
		s.RecipeID = recipe.ID
		err := s.Create(tx)
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
