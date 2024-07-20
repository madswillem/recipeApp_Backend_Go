package models

import (
	"fmt"
	"time"

	"github.com/madswillem/gompare"
	"github.com/madswillem/recipeApp_Backend_Go/internal/tools"
)

type RecipeGroupSchema struct {
	ID              string `db:"id"`
	IngredientDict  map[string]int
	IngredientVec   []float64
	PreperationDict map[string]int
	PreperationVec  []float64
	CuisineDict     map[string]int
	CuisineVec      []float64
	PrepTime        time.Duration
	CookingTime     time.Duration
	DietVec         []float64
	TechniquesDict  map[string]int
	TechniquesVec   []float64
	Recipes         []RecipeSchema
}

func (rp *RecipeGroupSchema) Create(r *RecipeSchema) {
	// Vectorize the Recipes
	// Vectorize Ingredients
	ingredient_list := make([]string, len(r.Ingredients))
	for n, i := range r.Ingredients {
		ingredient_list[n] = i.Name
	}

	h := gompare.New(gompare.Config{})
	h.Add(ingredient_list...)
	h.NormalMatrix()
	rp.IngredientDict = h.OuputMatrix.Dict
	rp.IngredientVec = tools.AverageVectors(h.OuputMatrix.Vec...)

	rp.Recipes = append(rp.Recipes, *r)

	h = gompare.New(gompare.Config{})
	h.Add(r.Cuisine)
	h.NormalMatrix()
	rp.CuisineDict = h.OuputMatrix.Dict
	rp.CuisineVec = h.OuputMatrix.Vec[0]

	technique_list := make([]string, len(r.Ingredients))
	for n, i := range r.Steps {
		technique_list[n] = *i.TechniqueID
	}
	h = gompare.New(gompare.Config{})
	h.Add(technique_list...)
	h.NormalMatrix()
	rp.TechniquesDict = h.OuputMatrix.Dict
	rp.TechniquesVec = tools.AverageVectors(h.OuputMatrix.Vec...)

	var hour, min, sec int
	fmt.Sscanf(r.PrepTime, "%d:%d:%d", hour, min, sec)
	rp.PrepTime = time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute + time.Duration(sec)*time.Second
	fmt.Sscanf(r.CookingTime, "%d:%d:%d", hour, min, sec)
	rp.CookingTime = time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute + time.Duration(sec)*time.Second
}
