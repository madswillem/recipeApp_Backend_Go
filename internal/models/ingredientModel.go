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