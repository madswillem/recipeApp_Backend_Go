// tests/s_test/CRUD_test.go
package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
	"github.com/madswillem/recipeApp_Backend_Go/internal/server"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func assertRecipesEqual(t *testing.T, expected models.RecipeSchema, actual models.RecipeSchema) {
	if actual.ID == 0 || expected.ID == 0 {
		t.Errorf("Expected non-nil values for actual and expected")
		return
	}

	// Check if the lengths of actual.Ingredients and expected.Ingredients are equal
	if len(actual.Ingredients) != len(expected.Ingredients) {
		t.Errorf("Expected %d ingredients but got %d", len(expected.Ingredients), len(actual.Ingredients))
		return
	}

	var errors []string

	// Compare each ingredient in the actual recipe
	for num, ingredient := range actual.Ingredients {
		expectedIngredient := expected.Ingredients[num]

		// Compare ingredient properties
		if ingredient.Ingredient != expectedIngredient.Ingredient {
			errors = append(errors, fmt.Sprintf("Expected ingredient %s but got %s", expectedIngredient.Ingredient, ingredient.Ingredient))
		}
		if ingredient.Amount != expectedIngredient.Amount {
			errors = append(errors, fmt.Sprintf("Expected amount %s but got %s", expectedIngredient.Amount, ingredient.Amount))
		}
		if ingredient.MeasurementUnit != expectedIngredient.MeasurementUnit {
			errors = append(errors, fmt.Sprintf("Expected measurement_unit %s but got %s", expectedIngredient.MeasurementUnit, ingredient.MeasurementUnit))
		}
		if ingredient.NutritionalValue != expectedIngredient.NutritionalValue {
			errors = append(errors, fmt.Sprintf("Expected nutritional_value %v but got %v", expectedIngredient.NutritionalValue, ingredient.NutritionalValue))
		}
		if ingredient.Rating != expectedIngredient.Rating {
			errors = append(errors, fmt.Sprintf("Expected rating %v but got %v", expectedIngredient.Rating, ingredient.Rating))
		}
	}

	// Compare other recipe properties
	if actual.Title != expected.Title {
		errors = append(errors, fmt.Sprintf("Expected title %s but got %s", expected.Title, actual.Title))
	}
	if actual.Preparation != expected.Preparation {
		errors = append(errors, fmt.Sprintf("Expected preparation %s but got %s", expected.Preparation, actual.Preparation))
	}
	if actual.CookingTime != expected.CookingTime {
		errors = append(errors, fmt.Sprintf("Expected cooking_time %d but got %d", expected.CookingTime, actual.CookingTime))
	}
	if actual.Image != expected.Image {
		errors = append(errors, fmt.Sprintf("Expected image %s but got %s", expected.Image, actual.Image))
	}
	if actual.NutriScore != expected.NutriScore {
		errors = append(errors, fmt.Sprintf("Expected nutriscore %s but got %s", expected.NutriScore, actual.NutriScore))
	}
	if actual.Rating != expected.Rating {
		errors = append(errors, fmt.Sprintf("Expected rating %v but got %v", expected.Rating, actual.Rating))
	}

	if len(errors) > 0 {
		t.Error(strings.Join(errors, "\n"))
	}
}
func readFileAsString(filePath string, t *testing.T) string {
	println(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Error opening file %s: %s", filePath, err)
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		t.Fatalf("Error copying file content: %s", err)
	}

	return buf.String()
}

func innitTestDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open("host=127.0.0.1 user=mads password=1234 dbname=test  port=5432 sslmode=disable"), &gorm.Config{})
	migrate(db)


	if err != nil || db == nil {
		panic("Error ")
	}
	return db
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.RecipeSchema{})
	db.AutoMigrate(&models.RecipeGroupSchema{})
	db.AutoMigrate(&models.Avrg{})
	db.AutoMigrate(&models.IngredientsSchema{})
	db.AutoMigrate(&models.RatingStruct{})
	db.AutoMigrate(&models.NutritionalValue{})
	db.AutoMigrate(&models.DietSchema{})
	db.AutoMigrate(&models.IngredientDBSchema{})
}

func TestAddRecipe(t *testing.T) {
	s := server.Server{DB: innitTestDB()}

	type testCase struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}

	testCases := []testCase{
		{
			name:           "add recipe with all required fields (edit=true)",
			requestBody:    "./testdata/create/add_recipe_with_all_required_fields(edited=true).json",
			expectedStatus: http.StatusCreated,
			expectedBody:   "./testdata/create/add_recipe_with_all_required_fields(edited=true)_expected_return.json",
		},
		{
			name:           "add recipe with all required fields (edit=false)",
			requestBody:    "./testdata/create/add_recipe_with_all_required_fields(edited=false).json",
			expectedStatus: http.StatusCreated,
			expectedBody:   "./testdata/create/add_recipe_with_all_required_fields(edited=false)_expected_return.json",
		},
		{
			name:           "add recipe with exesive edit fields",
			requestBody:    "./testdata/create/add_recipe_with_exesive_edit_fields.json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "add recipe with missing edit fields",
			requestBody:    "./testdata/create/add_recipe_with_missing_edit_fields.json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "add recipe without required fields",
			requestBody:    `./testdata/create/add_recipe_without_required_fields.json`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			completeRequestFilePath, err := filepath.Abs(tc.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			c.Request = httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(readFileAsString(completeRequestFilePath, t)))
			c.Request.Header.Set("Content-Type", "application/json")

			s.AddRecipe(c)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d but got %d. \n Body: %s", tc.expectedStatus, w.Code, w.Body.String())
			}

			if tc.expectedBody != "" {
				var response models.RecipeSchema
				err = json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatal(err)
				}

				completeExpectedFilePath, err := filepath.Abs(tc.expectedBody)
				if err != nil {
					t.Fatal(err)
				}

				var expectedReturn models.RecipeSchema
				expectedBody := readFileAsString(completeExpectedFilePath, t)
				err = json.Unmarshal([]byte(expectedBody), &expectedReturn)
				if err != nil {
					t.Fatal(err)
				}

				assertRecipesEqual(t, expectedReturn, response)
			}
		})
	}
}
func TestGetAll(t *testing.T) {
	s := server.Server{DB: innitTestDB()}
	t.Run("get all recipes", func(t *testing.T) {
		// Initialize the gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodGet, "/get", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the GetAll function
		s.GetAll(c)

		// Assert the response status code
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}

		// Assert the response body
		var response []models.RecipeSchema
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response body: %s", err.Error())
		}

		file_path, err := filepath.Abs("testdata/get/getAll.json")
		if err != nil {
			t.Fatal(err)
		}

		var expected_return []models.RecipeSchema
		err = json.Unmarshal([]byte(readFileAsString(file_path, t)), &expected_return)
		if err != nil {
			t.Errorf("Failed to unmarshal expected body: %s", err.Error())
		}

		for num, recipe := range response {
			assertRecipesEqual(t, expected_return[num], recipe)
		}
	})
}
func TestGetByID(t *testing.T) {
	s := server.Server{DB: innitTestDB()}
	
	t.Run("get recipe by id", func(t *testing.T) {
		// Initialize the gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		var err error
		c.Request, err = http.NewRequest(http.MethodGet, "/getbyid", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		if err != nil {
			t.Fatal(err)
		}

		// Call the GetById function
		s.GetById(c)

		// Assert the response status code
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d. \n Body: %s", http.StatusOK, w.Code, w.Body.String())
		}

		// Assert the response body
		var response models.RecipeSchema
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response body: %s", err.Error())
		}

		file_path, err := filepath.Abs("testdata/get/getByID.json")
		if err != nil {
			t.Fatal(err)
		}

		var expected_return models.RecipeSchema
		err = json.Unmarshal([]byte(readFileAsString(file_path, t)), &expected_return)
		if err != nil {
			t.Errorf("Failed to unmarshal expected body: %s", err.Error())
		}

		assertRecipesEqual(t, expected_return, response)
	})
	t.Run("get recipe by wrong id", func(t *testing.T) {
		// Initialize the gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodGet, "/recipes", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "200"}}

		// Call the GetById function
		s.GetById(c)

		// Assert the response status code
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d but got %d", http.StatusNotFound, w.Code)
		}
	})
}
func TestSelect(t *testing.T) {
	s := server.Server{DB: innitTestDB()}
	t.Run("select recipe", func(t *testing.T) {
		// Initialize the gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodGet, "/select", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		// Call the Select function
		s.Select(c)

		// Assert the response status code
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
			fmt.Println(w.Body.String())
		}
	})
	t.Run("select wrong recipe", func(t *testing.T) {
		// Initialize the gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodGet, "/select", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "200"}}

		// Call the Select function
		s.Select(c)

		// Assert the response status code
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d but got %d", http.StatusNotFound, w.Code)
		}

		// Assert the response body
		var body map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &body)
		if err != nil {
			t.Errorf("Failed to unmarshal response body: %s", err.Error())
		}

		if body["error"] == "error record not found" {
			t.Errorf("Expected message %s but got %s", "error record not found", body["error"])
		}
		if body["errMessage"] == "Recipe not found" {
			t.Errorf("Expected message %s but got %s", "Recipe not found", body["errorMessage"])

		}
	})
}
func TestDeselect(t *testing.T) {
	s := server.Server{DB: innitTestDB()}
	t.Run("deselect recipe", func(t *testing.T) {
		// Initialize the gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodGet, "/deselect", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		// Call the Select function
		s.Select(c)

		// Assert the response status code
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}
	})
	t.Run("deselect wrong recipe", func(t *testing.T) {
		// Initialize the gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodGet, "/deselect", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "200"}}

		// Call the Select function
		s.Select(c)

		// Assert the response status code
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d but got %d", http.StatusNotFound, w.Code)
		}

		// Assert the response body
		var body map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &body)
		if err != nil {
			t.Errorf("Failed to unmarshal response body: %s", err.Error())
		}

		if body["error"] == "error record not found" {
			t.Errorf("Expected message %s but got %s", "error record not found", body["error"])
		}
		if body["errMessage"] == "Recipe not found" {
			t.Errorf("Expected message %s but got %s", "Recipe not found", body["errorMessage"])

		}
	})
}
func TestDeleteRecipe(t *testing.T) {
	s := server.Server{DB: innitTestDB()}
	t.Run("delete recipe", func(t *testing.T) {
		// Initialize the gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodDelete, "/delete", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

		// Call the Delete function
		s.DeleteRecipe(c)

		// Assert the response status code
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
			fmt.Println(w.Body.String())
		}
	})
	t.Run("delete wrong recipe", func(t *testing.T) {
		// Initialize the gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodDelete, "/delete", nil)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "200"}}

		// Call the Delete function
		s.DeleteRecipe(c)

		// Assert the response status code
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d but got %d", http.StatusNotFound, w.Code)
		}
	})
}
func TestUpdateRecipe(t *testing.T) {
	s := server.Server{DB: innitTestDB()}
	t.Run("update recipe", func(t *testing.T) {
		// Initialize the gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Create a new request with the updated recipe
		updatedRecipe := make(map[string]interface{})
		updatedRecipe["title"] = "Updated Recipe"
		jsonBody, err := json.Marshal(updatedRecipe)
		if err != nil {
			t.Fatal(err)
		}

		c.Request, _ = http.NewRequest(http.MethodPatch, "/update", bytes.NewReader(jsonBody))
		c.Params = gin.Params{gin.Param{Key: "id", Value: "2"}}

		// Call the Update function
		s.UpdateRecipe(c)

		// Assert the response status code
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
			fmt.Println(w.Body.String())
		}
	})
	t.Run("delete wrong recipe", func(t *testing.T) {
		// Initialize the gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Create a new request with the updated recipe
		updatedRecipe := make(map[string]interface{})
		updatedRecipe["title"] = "Updated Wrong Recipe"
		jsonBody, err := json.Marshal(updatedRecipe)
		if err != nil {
			t.Fatal(err)
		}

		c.Request, _ = http.NewRequest(http.MethodPatch, "/update", bytes.NewReader(jsonBody))
		c.Params = gin.Params{gin.Param{Key: "id", Value: "5"}}

		// Call the Update function
		s.UpdateRecipe(c)

		// Assert the response status code
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d but got %d", http.StatusNotFound, w.Code)
			fmt.Println(w.Body.String())
		}
	})
}
