package models

type RecipeSchema struct {
	ID        		uint 				`json:"_id" gorm:"primarykey"`
	Title			string				`json:"title"` 
	Ingredients		[]IngredientsSchema	`json:"ingredients" gorm:"many2many:recipes_ingredients;"`
	Preparation		string				`json:"preparation"`
	Selected		int					`json:"selected"`	
	Version     	int					`json:"__v"`
}

type IngredientsSchema struct {
	ID        		uint 				`json:"id" gorm:"primarykey"`
	Ingredient		string				`json:"ingredient"` 
	Amount			string				`json:"amount"` 
}