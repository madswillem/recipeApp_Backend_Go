package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/madswillem/gompare"
	"github.com/madswillem/recipeApp_Backend_Go/internal/tools"
)

type RecipeGroupSchema struct {
	ID              string
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

func (rp *RecipeGroupSchema) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &rp)
		return nil
	case string:
		json.Unmarshal([]byte(v), &rp)
		return nil
	default:
		return fmt.Errorf("unsuported type: %T", v)
	}
}

func (rp *RecipeGroupSchema) Value() (driver.Value, error) {
	return json.Marshal(rp)
}

func (rp *RecipeGroupSchema) Create(r *RecipeSchema) {
	// Vectorize the Recipes
	// Vectorize Ingredients
	ingredient_list := make([]string, len(r.Ingredients))
	for n, i := range r.Ingredients {
		ingredient_list[n] = i.Name
	}

	h := gompare.New(gompare.Config{})
	h.InputStrings = make([][]string, 1)
	h.InputStrings[0] = ingredient_list
	h.NormalMatrix()
	rp.IngredientDict = h.OutputMatrix.Dict
	rp.IngredientVec = tools.AverageVectors(h.OutputMatrix.Vec...)

	rp.Recipes = append(rp.Recipes, *r)

	h = gompare.New(gompare.Config{})
	h.Add(r.Cuisine)
	h.NormalMatrix()
	rp.CuisineDict = h.OutputMatrix.Dict
	rp.CuisineVec = h.OutputMatrix.Vec[0]

	technique_list := make([]string, len(r.Ingredients))
	for n, i := range r.Steps {
		if i.TechniqueID == nil {
			t := ""
			i.TechniqueID = &t
		}
		technique_list[n] = *i.TechniqueID
	}
	h = gompare.New(gompare.Config{})
	h.InputStrings = append(h.InputStrings, technique_list)
	h.NormalMatrix()
	rp.TechniquesDict = h.OutputMatrix.Dict
	rp.TechniquesVec = tools.AverageVectors(h.OutputMatrix.Vec...)

	var hour, min, sec int
	fmt.Sscanf(r.PrepTime, "%d:%d:%d", &hour, &min, &sec)
	rp.PrepTime = time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute + time.Duration(sec)*time.Second
	fmt.Sscanf(r.CookingTime, "%d:%d:%d", hour, min, sec)
	rp.CookingTime = time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute + time.Duration(sec)*time.Second

	steps := make([]string, len(r.Steps))
	for i, s := range r.Steps {
		steps[i] = s.Step
	}
	prep := strings.Join(steps, " ")
	h = gompare.New(gompare.Config{})
	h.Add(prep)
	h.NormalMatrix()
	rp.PreperationDict = h.OutputMatrix.Dict
	rp.PreperationVec = h.OutputMatrix.Vec[0]
}

func (rp *RecipeGroupSchema) Compare(r *RecipeSchema) float64 {
	h := gompare.New(gompare.Config{
		Matrix: gompare.Matrix{
			Dict: rp.IngredientDict,
			Vec:  [][]float64{rp.IngredientVec},
		},
	})
	ingredient_list := make([]string, len(r.Ingredients))
	for n, i := range r.Ingredients {
		ingredient_list[n] = i.Name
	}
	h.InputStrings = append(h.InputStrings, ingredient_list)
	h.NormalMatrix()
	h.CosineSimilarity(0, 1)
	simIngs := h.Similarity

	h = gompare.New(gompare.Config{
		Matrix: gompare.Matrix{
			Dict: rp.CuisineDict,
			Vec:  [][]float64{rp.CuisineVec},
		},
	})
	h.Add(r.Cuisine)
	h.NormalMatrix()
	h.CosineSimilarity(0, 1)
	simCuisine := h.Similarity

	h = gompare.New(gompare.Config{
		Matrix: gompare.Matrix{
			Dict: rp.PreperationDict,
			Vec:  [][]float64{rp.PreperationVec},
		},
	})
	steps := make([]string, len(r.Steps))
	for i, s := range r.Steps {
		steps[i] = s.Step
	}
	prep := strings.Join(steps, " ")
	h.Add(prep)
	h.NormalMatrix()
	h.CosineSimilarity(0, 1)
	simPrep := h.Similarity

	technique_list := make([]string, len(r.Ingredients))
	for n, i := range r.Steps {
		if i.TechniqueID == nil {
			t := ""
			i.TechniqueID = &t
		}
		technique_list[n] = *i.TechniqueID
	}
	h = gompare.New(gompare.Config{
		Matrix: gompare.Matrix{
			Dict: rp.TechniquesDict,
			Vec:  [][]float64{rp.TechniquesVec},
		},
	})
	h.InputStrings = append(h.InputStrings, technique_list)
	h.NormalMatrix()
	h.CosineSimilarity(0, 1)
	simTech := h.Similarity

	sim := simIngs*2 + simCuisine*3 + simPrep + simTech*2
	return sim / 8
}

func (rp *RecipeGroupSchema) Add(r *RecipeSchema) {
	ingredient_list := make([]string, len(r.Ingredients))
	for n, i := range r.Ingredients {
		ingredient_list[n] = i.Name
	}

	h := gompare.New(gompare.Config{
		Matrix: gompare.Matrix{
			Dict: rp.IngredientDict,
		},
	})
	h.InputStrings = make([][]string, 1)
	h.InputStrings[0] = ingredient_list
	h.NormalMatrix()
	vec := tools.AddVectors(tools.MultiplyVectorByNum(float64(len(rp.Recipes)), rp.IngredientVec), h.OutputMatrix.Vec[0])
	rp.IngredientDict = h.OutputMatrix.Dict
	rp.IngredientVec = tools.MultiplyVectorByNum(1.0/float64(len(rp.Recipes)+1), vec)

	steps := make([]string, len(r.Steps))
	for i, s := range r.Steps {
		steps[i] = s.Step
	}
	prep := strings.Join(steps, " ")
	h = gompare.New(gompare.Config{
		Matrix: gompare.Matrix{
			Dict: rp.PreperationDict,
		},
	})
	h.Add(prep)
	h.NormalMatrix()
	vec = tools.AddVectors(tools.MultiplyVectorByNum(float64(len(rp.Recipes)), rp.PreperationVec), h.OutputMatrix.Vec[0])
	rp.PreperationDict = h.OutputMatrix.Dict
	rp.PreperationVec = tools.MultiplyVectorByNum(1.0/float64(len(rp.Recipes)+1), vec)

	h = gompare.New(gompare.Config{
		Matrix: gompare.Matrix{
			Dict: rp.CuisineDict,
		},
	})
	h.Add(r.Cuisine)
	h.NormalMatrix()
	vec = tools.AddVectors(tools.MultiplyVectorByNum(float64(len(rp.Recipes)), rp.CuisineVec), h.OutputMatrix.Vec[0])
	rp.CuisineDict = h.OutputMatrix.Dict
	rp.CuisineVec = tools.MultiplyVectorByNum(1.0/float64(len(rp.Recipes)+1), vec)

	var hour, min, sec int
	fmt.Sscanf(r.PrepTime, "%d:%d:%d", &hour, &min, &sec)
	rp.PrepTime += time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute + time.Duration(sec)*time.Second
	fmt.Println(time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute + time.Duration(sec)*time.Second)
	fmt.Println(len(rp.Recipes) + 1)
	rp.PrepTime /= time.Duration(len(rp.Recipes) + 1)
	fmt.Sscanf(r.CookingTime, "%d:%d:%d", &hour, &min, &sec)
	rp.CookingTime += time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute + time.Duration(sec)*time.Second
	rp.CookingTime /= time.Duration(len(rp.Recipes) + 1)

	return
}
