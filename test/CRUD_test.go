package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/madswillem/recipeApp_Backend_Go/internal/database"
	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
	"github.com/madswillem/recipeApp_Backend_Go/internal/server"
	"github.com/madswillem/recipeApp_Backend_Go/internal/tools"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func assertRecipesEqual(t *testing.T, expected models.RecipeSchema, actual models.RecipeSchema) {
	if actual.ID == "" || expected.ID == "" {
		t.Errorf(fmt.Sprintf("%s (actual.ID) != %s (expected.ID)", actual.ID, expected.ID))

		return
	}

	// Check if the lengths of actual.Ingredients and expected.Ingredients are equal
	if len(actual.Ingredients) != len(expected.Ingredients) {
		t.Errorf("Expected %d ingredients but got %d", len(expected.Ingredients), len(actual.Ingredients))
		return
	}

	// Check if the lengths of actual.Steps and expected.Steps are equal
	if len(actual.Steps) != len(expected.Steps) {
		t.Errorf("Expected %d ingredients but got %d", len(expected.Steps), len(actual.Steps))
		return
	}

	var errors []string

	// Compare each ingredient in the actual recipe
	for num, ingredient := range actual.Ingredients {
		expectedIngredient := expected.Ingredients[num]

		// Compare ingredient properties
		if ingredient.Name != expectedIngredient.Name {
			errors = append(errors, fmt.Sprintf("Expected ingredient %s but got %s", expectedIngredient.Name, ingredient.Name))
		}
		if ingredient.Amount != expectedIngredient.Amount {
			errors = append(errors, fmt.Sprintf("Expected amount %d but got %d", expectedIngredient.Amount, ingredient.Amount))
		}
		if ingredient.Unit != expectedIngredient.Unit {
			errors = append(errors, fmt.Sprintf("Expected measurement_unit %s but got %s", expectedIngredient.Unit, ingredient.Unit))
		}
		if ingredient.NutritionalValue != expectedIngredient.NutritionalValue {
			errors = append(errors, fmt.Sprintf("Expected nutritional_value %v but got %v", expectedIngredient.NutritionalValue, ingredient.NutritionalValue))
		}
		if ingredient.Rating != expectedIngredient.Rating {
			errors = append(errors, fmt.Sprintf("Expected rating %v but got %v", expectedIngredient.Rating, ingredient.Rating))
		}
	}

	// Compare Steps
	for num, step := range actual.Steps {
		expectedStep := expected.Steps[num]

		if step.Step != expectedStep.Step {
			errors = append(errors, fmt.Sprintf("Expected step %s but got %s", expectedStep.Step, step.Step))
		}
		if step.TechniqueID != expectedStep.TechniqueID {
			errors = append(errors, fmt.Sprintf("Expected technique %s but got %s", *expectedStep.TechniqueID, *step.TechniqueID))
		}
	}

	// Compare other recipe properties
	if actual.Name != expected.Name {
		errors = append(errors, fmt.Sprintf("Expected Name %s but got %s", expected.Name, actual.Name))
	}
	if actual.PrepTime != expected.PrepTime {
		errors = append(errors, fmt.Sprintf("Expected prep_time %s but got %s", expected.PrepTime, actual.PrepTime))
	}
	if actual.CookingTime != expected.CookingTime {
		errors = append(errors, fmt.Sprintf("Expected cooking_time %s but got %s", expected.CookingTime, actual.CookingTime))
	}
	if actual.NutritionalValue != expected.NutritionalValue {
		errors = append(errors, fmt.Sprintf("Expected NutritionalValue %v but got %v", expected.NutritionalValue, actual.NutritionalValue))
	}
	if actual.Rating != expected.Rating {
		errors = append(errors, fmt.Sprintf("Expected rating %v but got %v", expected.Rating, actual.Rating))
	}

	if len(errors) > 0 {
		t.Error(strings.Join(errors, "\n"))
	}
}

func TestServer_AddRecipe(t *testing.T) {
	ctx := context.Background()

	// 1. Start the postgres container and run any migrations on it
	container, err := postgres.Run(
		ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("mads"),
		postgres.WithPassword("1234"),
		postgres.BasicWaitStrategies(),
		postgres.WithInitScripts("./testdata/innit-db.sql"),
		postgres.WithSQLDriver("pgx"),
	)
	if err != nil {
		t.Fatal(err)
	}

	// 2. Create a snapshot of the database to restore later
	err = container.Snapshot(ctx, postgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	dbURL, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	s := server.Server{NewDB: database.ConnectToDB(&sqlx.Conn{}, dbURL)}

	type testCase struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}

	testCases := []testCase{
		{
			name:           "add recipe with all required fields",
			requestBody:    "./testdata/create/add_recipe_with_all_required_fields/body.json",
			expectedStatus: http.StatusCreated,
			expectedBody:   "./testdata/create/add_recipe_with_all_required_fields/expected_return.json",
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

			reqBody, err := tools.ReadFileAsString(completeRequestFilePath)
			if err != nil {
				t.Error(err)
				return
			}
			c.Request = httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(reqBody))
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
				expectedBody, err := tools.ReadFileAsString(completeExpectedFilePath)
				if err != nil {
					t.Fatal(err)
				}
				err = json.Unmarshal([]byte(expectedBody), &expectedReturn)
				if err != nil {
					t.Fatal(err)
				}

				assertRecipesEqual(t, expectedReturn, response)
			}

		})
	}
}
