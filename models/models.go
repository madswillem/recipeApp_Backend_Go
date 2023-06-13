package models

type RecipeSchema struct {
	ID        		string 				`json:"_id" gorm:"primarykey"`
	Title			string				`json:"title" gorm:"unique"` 
	Ingredients		[]IngredientsSchema	`json:"ingredients" gorm:"foreignKey:RecipeSchemaID"`
	Preparation		string				`json:"preparation"`
	CookingTime		int					`json:"cookingtime"`
	Selected		int					`json:"selected"`
	Rating			RatingStruct		`json:"rating" gorm:"polymorphic:Owner"`	
	Version     	int					`json:"__v"`
}

type IngredientsSchema struct {
	ID        		string 				`json:"id" gorm:"primarykey"`
	Ingredient		string				`json:"ingredient"` 
	Amount			string				`json:"amount"` 
	Rating			RatingStruct		`json:"rating" gorm:"polymorphic:Owner"`
	RecipeSchemaID 	string				`json:"-"`
}

type CurrentData struct {
	Day				string
	Season			string
	Temp			string
}