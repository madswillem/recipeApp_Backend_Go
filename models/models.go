package models

type RecipeSchema struct {
	ID        		 string 			 `json:"_id" gorm:"primarykey"`
	Title			 string				 `json:"title" gorm:"unique"` 
	Ingredients		 []IngredientsSchema `json:"ingredients" gorm:"foreignKey:RecipeSchemaID"`
	Preparation		 string				 `json:"preparation"`
	CookingTime		 int				 `json:"cookingtime"`
	Image			 string				 `json:"image"`
	NutriScore		 string				 `json:"nutriscore"`
	NutritionalValue NutritionalValue    `json:"nutritional_value" gorm:"polymorphic:Owner"`
	Diet			 DietSchema			 `json:"diet" gorm:"foreignKey:RecipeID"`
	Selected		 int				 `json:"selected"`
	Rating			 RatingStruct		 `json:"rating" gorm:"polymorphic:Owner"`	
	Version     	 int				 `json:"__v"`
}

type IngredientsSchema struct {
	ID        		 string 			 `json:"id" gorm:"primarykey"`
	Ingredient		 string				 `json:"ingredient"` 
	Amount			 string				 `json:"amount"` 
	NutritionalValue NutritionalValue    `json:"nutritional_value" gorm:"polymorphic:Owner"`
	Rating			 RatingStruct		 `json:"rating" gorm:"polymorphic:Owner"`
	RecipeSchemaID 	 string				 `json:"-"`
}

type DietSchema struct {
	ID		  		 string 			 `gorm:"primarykey"`
	RecipeID 	     string				 

	Vegetarien		 bool				 `json:"vegetarien"`
	Vegan			 bool				 `json:"vegan"`
	LowCal			 bool				 `json:"lowcal"`
	LowCarb			 bool				 `json:"lowcarb"`
	Keto			 bool				 `json:"keto"`
	Paleo			 bool				 `json:"paleo"`
	LowFat			 bool				 `json:"lowfat"`
	FoodCombining	 bool				 `json:"food_combining"`
	WholeFood		 bool				 `json:"whole_food"`
}