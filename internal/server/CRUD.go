package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
)

func (s *Server) GetAll(c *gin.Context) {
	recipes := []models.RecipeSchema{}
	ingredients := []models.IngredientsSchema{}
	steps := []models.StepsStruct{}
	nutritional_values := []models.NutritionalValue{}

	err := s.NewDB.Select(&recipes, "SELECT * FROM recipes")
	if err != nil {
		print(err.Error())
		error_handler.HandleError(c, http.StatusBadRequest, "Error while getting recipes", []error{err})
		return
	}

	err = s.NewDB.Select(&ingredients, "SELECT ri.*, i.name AS name FROM recipe_ingredient ri JOIN ingredient i ON ri.ingredient_id = i.id")
	if err != nil {
		print(err.Error())
		error_handler.HandleError(c, http.StatusBadRequest, "Error while getting ingredients", []error{err})
		return
	}

	err = s.NewDB.Select(&steps, "SELECT * FROM step")
	if err != nil {
		print(err.Error())
		error_handler.HandleError(c, http.StatusBadRequest, "Error while getting steps", []error{err})
		return
	}

	err = s.NewDB.Select(&nutritional_values, "SELECT * FROM nutritional_value WHERE (recipe_id IS NOT NULL)::integer = 1")
	if err != nil {
		print(err.Error())
		error_handler.HandleError(c, http.StatusBadRequest, "Error while getting nutritional values", []error{err})
		return
	}

	recipeMap := make(map[string]*models.RecipeSchema)
	for i := range recipes {
		recipeMap[recipes[i].ID] = &recipes[i]
	}

	for _, ingredient := range ingredients {
		if recipe, found := recipeMap[ingredient.RecipeID]; found {
			recipe.Ingredients = append(recipe.Ingredients, ingredient)
		}
	}
	for _, step := range steps {
		if recipe, found := recipeMap[*step.RecipeID]; found {
			recipe.Steps = append(recipe.Steps, step)
		}
	}
	c.JSON(http.StatusOK, recipes)
}

func (s *Server) AddRecipe(c *gin.Context) {
	var body models.RecipeSchema

	binderr := c.ShouldBindJSON(&body)
	if binderr != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "Failed to read body", []error{binderr})
		return
	}

	err := body.Create(s.NewDB)
	if err != nil {
		error_handler.HandleError(c, err.Code, err.Message, err.Errors)
		return
	}

	c.JSON(http.StatusCreated, body)
}

func (s *Server) AddIngredient(c *gin.Context) {
	var body models.IngredientDB

	binderr := c.ShouldBindJSON(&body)
	if binderr != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "Failed to read body", []error{binderr})
		return
	}

	err := body.Create(s.NewDB)
	if err != nil {
		error_handler.HandleError(c, err.Code, err.Message, err.Errors)
		return
	}

	c.JSON(http.StatusAccepted, body)
}

func (s *Server) UpdateRecipe(c *gin.Context) {

	var body struct {
		Name        *string `db:"name" json:"name"`
		Cuisine     *string `db:"cuisine" json:"cuisine"`
		Yield       *int    `db:"yield" json:"yield"`
		YieldUnit   *string `db:"yield_unit" json:"yield_unit"`
		PrepTime    *string `db:"prep_time" json:"prep_time"`
		CookingTime *string `db:"cooking_time" json:"cooking_time"`
		Ingredients *[]models.IngredientsSchema
		Diet        *models.DietSchema
		Steps       *[]models.StepsStruct
	}

	i := c.Param("id")

	c.ShouldBindJSON(&body)

	tx, err := s.NewDB.Beginx()
	if err != nil {
		error_handler.HandleError(c, http.StatusInternalServerError, "Error initiating transaction", []error{err})
	}

	var setParts []string
	var args []interface{}

	if body.Name != nil {
		setParts = append(setParts, "name = $"+strconv.Itoa(len(args)+1))
		args = append(args, *body.Name)
	}
	if body.Cuisine != nil {
		setParts = append(setParts, "cuisine = $"+strconv.Itoa(len(args)+1))
		args = append(args, *body.Cuisine)
	}
	if body.Yield != nil {
		setParts = append(setParts, "yield = $"+strconv.Itoa(len(args)+1))
		args = append(args, *body.Yield)
	}
	if body.YieldUnit != nil {
		setParts = append(setParts, "yield_unit = $"+strconv.Itoa(len(args)+1))
		args = append(args, *body.YieldUnit)
	}
	if body.PrepTime != nil {
		setParts = append(setParts, "prep_time = $"+strconv.Itoa(len(args)+1))
		args = append(args, *body.PrepTime)
	}
	if body.CookingTime != nil {
		setParts = append(setParts, "cooking_time = $"+strconv.Itoa(len(args)+1))
		args = append(args, *body.CookingTime)
	}

	if len(setParts) == 0 {
		// No fields to body
		error_handler.HandleError(c, http.StatusExpectationFailed, "Nothing to update", []error{})
		return
	}

	query := "UPDATE recipes SET " + strings.Join(setParts, ", ") + " WHERE id = $" + strconv.Itoa(len(args)+1)
	fmt.Println(query)
	args = append(args, i)

	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		error_handler.HandleError(c, http.StatusInternalServerError, "Error Updating recipe", []error{err})
		return
	}

	tx.Commit()
}

func (s *Server) DeleteRecipe(c *gin.Context) {
	i := c.Param("id")
	err := s.NewDB.QueryRowx(`DELETE FROM public.recipes WHERE id = $1`, i).Err()
	if err != nil {
		error_handler.HandleError(c, http.StatusInternalServerError, err.Error(), []error{err})
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) GetById(c *gin.Context) {
	i := c.Param("id")

	result := models.RecipeSchema{ID: i}
	err := result.GetRecipeByID(s.NewDB)

	if err != nil {
		error_handler.HandleError(c, err.Code, err.Message, err.Errors)
	}
	c.JSON(http.StatusOK, result)
}

func (s *Server) Filter(c *gin.Context) {

	var body models.Filter
	err := c.ShouldBindJSON(&body)

	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "Failed to read body", []error{err})
		return
	}

	recipes, apiErr := body.Filter(s.NewDB)
	if apiErr != nil {
		error_handler.HandleError(c, apiErr.Code, apiErr.Message, apiErr.Errors)
		return
	}

	c.JSON(http.StatusOK, recipes)
}

func (s *Server) Select(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "id is not a number", []error{err})
		return
	}
	middleware_user, _ := c.Get("user")
	user, ok := middleware_user.(models.UserModel)
	if !ok {
		fmt.Println("type assertion failed")
	}

	response := models.RecipeSchema{}
	response.ID = fmt.Sprint(i)

	selectedErr := response.UpdateSelected(1, &user, s.DB)
	if selectedErr != nil {
		error_handler.HandleError(c, selectedErr.Code, selectedErr.Message, selectedErr.Errors)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) Deselect(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_handler.HandleError(c, http.StatusBadRequest, "id is not a number", []error{err})
		return
	}
	middleware_user, _ := c.MustGet("user").(models.UserModel)

	response := models.RecipeSchema{}
	response.ID = fmt.Sprint(i)

	selectedErr := response.UpdateSelected(-1, &middleware_user, s.DB)
	if selectedErr != nil {
		error_handler.HandleError(c, selectedErr.Code, selectedErr.Message, selectedErr.Errors)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) Colormode(c *gin.Context) {
	switch c.Param("type") {
	case "get":
		cookie, err := c.Cookie("type")
		if err != nil {
			error_handler.HandleError(c, http.StatusBadRequest, "Cookie error", []error{err})
		}
		c.JSON(http.StatusOK, gin.H{"type": cookie})
	case "dark":
		c.SetCookie("type", "dark", 999999999999999999, "/", "localhost", false, true)
		c.Status(http.StatusAccepted)
	case "light":
		c.SetCookie("type", "light", 999999999999999999, "/", "localhost", false, true)
		c.Status(http.StatusAccepted)
	default:
		c.Status(http.StatusBadRequest)
	}
}
