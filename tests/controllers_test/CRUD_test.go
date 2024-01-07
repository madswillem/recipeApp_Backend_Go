package controllers_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"rezeptapp.ml/goApp/controllers"
	"rezeptapp.ml/goApp/models"
)

func readFileAsString(filePath string, t *testing.T) string {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Error reading file %s: %s", filePath, err)
	}
	return string(fileContent)
}

func TestAddRecipe(t *testing.T) {
	t.Run("add recipe with all required fields (edit=true)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
	
		// Set the request body
		requestBody := readFileAsString("../testdata/create/add_recipe_with_all_required_fields.json", t)
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
	
		// Assert the recipe title
		if response.Title != "Test Recipe" {
			t.Errorf("Expected recipe title 'Test Recipe' but got '%s'", response.Title)
		}
	
		// Assert the number of ingredients
		if len(response.Ingredients) != 2 {
			t.Errorf("Expected 2 ingredients but got %d", len(response.Ingredients))
		}
	
		// Assert the first ingredient
		if response.Ingredients[0].Ingredient != "Ingredient 1" {
			t.Errorf("Expected first ingredient 'Ingredient 1' but got '%s'", response.Ingredients[0].Ingredient)
		}
	
		// Assert the second ingredient
		if response.Ingredients[1].Ingredient != "Ingredient 2" {
			t.Errorf("Expected second ingredient 'Ingredient 2' but got '%s'", response.Ingredients[1].Ingredient)
		}
	
		// Assert the recipe instructions
		if response.Preparation!= "Step 1: Do this, Step 2: Do that" {
			t.Errorf("Expected recipe preparation 'Step 1: Do this, Step 2: Do that' but got '%s'", response.Preparation)
		}
	})
	t.Run("add recipe with all required fields (edit=false)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
	
		// Set the request body
		requestBody := readFileAsString("../testdata/create/add_recipe_with_all_required_fields.json", t)
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
	
		// Assert the recipe title
		if response.Title != "Test Recipe" {
			t.Errorf("Expected recipe title 'Test Recipe' but got '%s'", response.Title)
		}
	
		// Assert the number of ingredients
		if len(response.Ingredients) != 2 {
			t.Errorf("Expected 2 ingredients but got %d", len(response.Ingredients))
		}
	
		// Assert the first ingredient
		if response.Ingredients[0].Ingredient != "Ingredient 1" {
			t.Errorf("Expected first ingredient 'Ingredient 1' but got '%s'", response.Ingredients[0].Ingredient)
		}
	
		// Assert the second ingredient
		if response.Ingredients[1].Ingredient != "Ingredient 2" {
			t.Errorf("Expected second ingredient 'Ingredient 2' but got '%s'", response.Ingredients[1].Ingredient)
		}
	
		// Assert the recipe instructions
		if response.Preparation!= "Step 1: Do this, Step 2: Do that" {
			t.Errorf("Expected recipe preparation 'Step 1: Do this, Step 2: Do that' but got '%s'", response.Preparation)
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
	
		// Assert the response body
		var response models.RecipeSchema
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response body: %s", err.Error())
		}
	
		// Assert the recipe title
		if response.Title != "Test Recipe" {
			t.Errorf("Expected recipe title 'Test Recipe' but got '%s'", response.Title)
		}
	
		// Assert the number of ingredients
		if len(response.Ingredients) != 0 {
			t.Errorf("Expected 0 ingredients but got %d", len(response.Ingredients))
		}
	
		// Assert the recipe instructions
		if response.Preparation != "Step 1: Do this, Step 2: Do that" {
			t.Errorf("Expected recipe instructions 'Step 1: Do this, Step 2: Do that' but got '%s'", response.Preparation)
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
	
		// Assert the response body
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response body: %s", err.Error())
		}
	
		// Assert the error message
		errorMessage, ok := response["error"].(string)
		if !ok || errorMessage != "Failed to read body" {
			t.Errorf("Expected error message 'Failed to read body' but got '%s'", errorMessage)
		}
	})
}