package models

type RecipeSchema struct {
	ID        		string 				`json:"_id" gorm:"primarykey"`
	Title			string				`json:"title" gorm:"unique"` 
	Ingredients		[]*IngredientsSchema`json:"ingredients" gorm:"many2many:recipes_ingredients;"`
	Preparation		string				`json:"preparation"`
	Selected		int					`json:"selected"`	
	Version     	int					`json:"__v"`
}

type IngredientsSchema struct {
	ID        		string 				`json:"id" gorm:"primarykey"`
	Ingredient		string				`json:"ingredient" gorm:"unique"` 
	Amount			string				`json:"amount"` 
	Recipes			[]*RecipeSchema		`json:"-" gorm:"many2many:recipes_ingredients;"`
}