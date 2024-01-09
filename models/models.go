package models

type IngredientsSchema struct {
	ID        		 uint 			 	 `json:"id" gorm:"primarykey"`
	RecipeSchemaID 	 uint				 `json:"-"`

	Ingredient		 string				 `json:"ingredient"` 
	Amount			 string				 `json:"amount"` 
	MeasurementUnit	 string				 `json:"measurement_unit"`
	NutritionalValue NutritionalValue    `json:"nutritional_value" gorm:"polymorphic:Owner"`
	Rating			 RatingStruct		 `json:"rating" gorm:"polymorphic:Owner"`
}

type DietSchema struct {
	ID		  		 uint 			 	 `gorm:"primarykey"`
	OwnerTitle	string 					 `json:"owner_title"`
	OwnerID   	string
  	OwnerType 	string		 

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