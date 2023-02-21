package tools

import (
	"rezeptapp.ml/goApp/initializers"
)

func CheckIfRecipeExists(title string) (bool, error) {
	var result struct {
		Found bool
	}
	  
	err := initializers.DB.Raw("SELECT EXISTS(SELECT * FROM ingredients_schemas WHERE ingredient = ?) AS found;", title).Scan(&result).Error

	return result.Found, err
}