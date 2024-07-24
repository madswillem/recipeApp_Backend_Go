package test

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/madswillem/recipeApp_Backend_Go/internal/database"
	"github.com/madswillem/recipeApp_Backend_Go/internal/models"
)

func TestSimilarity(t *testing.T) {
	t.Run("test similarity recipe_group recipe", func(t *testing.T) {
		db := database.ConnectToDB(&sqlx.Conn{})
		r := models.RecipeSchema{ID: "aa85daf1-dbc5-462d-a6fe-3fbb358b08dd"}
		r.GetRecipeByID(db)
		rp := models.RecipeGroupSchema{}
		rp.Create(&r)

		fmt.Println(rp.Compare(&r))
	})
}
