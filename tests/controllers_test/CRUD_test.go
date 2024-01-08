// tests/controllers_test/CRUD_test.go
package controllers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"rezeptapp.ml/goApp/controllers"
	"rezeptapp.ml/goApp/models"
)

// assertRecipesEqual compares two RecipeSchema structs and checks if their values are equal.
// It is primarily used in testing to verify that the expected and actual recipe data match.
func assertRecipesEqual(t *testing.T, expected models.RecipeSchema, actual models.RecipeSchema) {
	// Compare each ingredient in the actual recipe
	for num, ingredient := range actual.Ingredients {
		expectedIngredient := expected.Ingredients[num]

		// Compare ingredient properties
		if ingredient.Ingredient != expectedIngredient.Ingredient {
			t.Errorf("Expected ingredient %s but got %s", expectedIngredient.Ingredient, ingredient.Ingredient)
		}
		if ingredient.Amount != expectedIngredient.Amount {
			t.Errorf("Expected amount %s but got %s", expectedIngredient.Amount, ingredient.Amount)
		}
		if ingredient.MeasurementUnit != expectedIngredient.MeasurementUnit {
			t.Errorf("Expected measurement_unit %s but got %s", expectedIngredient.MeasurementUnit, ingredient.MeasurementUnit)
		}
		if ingredient.NutritionalValue != expectedIngredient.NutritionalValue {
			t.Errorf("Expected nutritional_value %v but got %v", expectedIngredient.NutritionalValue, ingredient.NutritionalValue)
		}
		if ingredient.Rating != expectedIngredient.Rating {
			t.Errorf("Expected rating %v but got %v", expectedIngredient.Rating, ingredient.Rating)
		}
	}

	// Compare other recipe properties
	if actual.Title != expected.Title {
		t.Errorf("Expected title %s but got %s", expected.Title, actual.Title)
	}
	if actual.Preparation != expected.Preparation {
		t.Errorf("Expected preparation %s but got %s", expected.Preparation, actual.Preparation)
	}
	if actual.CookingTime != expected.CookingTime {
		t.Errorf("Expected cooking_time %d but got %d", expected.CookingTime, actual.CookingTime)
	}
	if actual.Image != expected.Image {
		t.Errorf("Expected image %s but got %s", expected.Image, actual.Image)
	}
	if actual.NutriScore != expected.NutriScore {
		t.Errorf("Expected nutriscore %s but got %s", expected.NutriScore, actual.NutriScore)
	}
	if actual.Rating != expected.Rating {
		t.Errorf("Expected rating %v but got %v", expected.Rating, actual.Rating)
	}
}
func readFileAsString(filePath string, t *testing.T) string {
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

func TestAddRecipe(t *testing.T) {
	t.Run("add recipe with all required fields (edit=true)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set the request body
		requestBody := readFileAsString("../testdata/create/add_recipe_with_all_required_fields(edit=true).json", t)
		c.Request, _ = http.NewRequest(http.MethodPost, "/recipes", strings.NewReader(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the AddRecipe function
		controllers.AddRecipe(c)

		// Assert the response status code
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d but got %d", http.StatusCreated, w.Code)
		}

		// Assert the response body
		var response models.RecipeSchema
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response body: %s", err.Error())
		}
		var expected_return models.RecipeSchema
		err = json.Unmarshal([]byte(readFileAsString("../testdata/create/add_recipe_with_all_required_fields_(edit=true)_expected_return.json", t)), &expected_return)
		if err != nil {
			t.Errorf("Failed to unmarshal response body: %s", err.Error())
		}

		assertRecipesEqual(t, expected_return, response)
	})
	t.Run("add recipe with all required fields (edit=false)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set the request body
		requestBody := readFileAsString("../testdata/create/add_recipe_with_all_required_fields(edit=false).json", t)
		c.Request, _ = http.NewRequest(http.MethodPost, "/recipes", strings.NewReader(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the AddRecipe function
		controllers.AddRecipe(c)

		// Assert the response status code
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d but got %d", http.StatusCreated, w.Code)
		}

		// Assert the response body
		var response models.RecipeSchema
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response body: %s", err.Error())
		}
		var expected_return models.RecipeSchema
		err = json.Unmarshal([]byte(readFileAsString("../testdata/create/add_recipe_with_all_required_fields_(edit=false)_expected_return.json", t)), &expected_return)
		if err != nil {
			t.Errorf("Failed to unmarshal response body: %s", err.Error())
		}

		assertRecipesEqual(t, expected_return, response)
	})
	t.Run("add recipe with exesive edit fields", func(t *testing.T) {
		// Initialize the gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set the request body with missing required fields
		requestBody := `{
			"description": "This is a test recipe",
			"ingredients": [
				{
					"ingredient": "Ingredient 5",
					"quantity": 2
				},
				{
					"ingredient": "Ingredient 6",
					"quantity": 3
				}
			],
			"instructions": "Step 1: Do this, Step 2: Do that"
		}`
		c.Request, _ = http.NewRequest(http.MethodPost, "/recipes", strings.NewReader(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the AddRecipe function
		controllers.AddRecipe(c)

		// Assert the response status code
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}
	})
	t.Run("add recipe with missing edit fields", func(t *testing.T) {
		// Initialize the gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set the request body with missing required fields
		requestBody := `{
			"description": "This is a test recipe",
			"ingredients": [
				{
					"ingredient": "Ingredient 5",
					"quantity": 2
				},
				{
					"ingredient": "Ingredient 6",
					"quantity": 3
				}
			],
			"instructions": "Step 1: Do this, Step 2: Do that"
		}`
		c.Request, _ = http.NewRequest(http.MethodPost, "/recipes", strings.NewReader(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the AddRecipe function
		controllers.AddRecipe(c)

		// Assert the response status code
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}
	})
	t.Run("add recipe without required fields", func(t *testing.T) {
		// Initialize the gin context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set the request body
		requestBody := `{
			"title": "Test Recipe",
			"Preperation": "Step 1: Do this, Step 2: Do that"
		}`
		c.Request, _ = http.NewRequest(http.MethodPost, "/recipes", strings.NewReader(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the AddRecipe function
		controllers.AddRecipe(c)

		// Assert the response status code
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d but got %d", http.StatusCreated, w.Code)
		}
	})
}
