package test

import (
	"context"
	"encoding/json"
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

	println("test")
	for _, tc := range testCases {
		println("test")
		t.Run(tc.name, func(t *testing.T) {
			println("test " + tc.name)
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
