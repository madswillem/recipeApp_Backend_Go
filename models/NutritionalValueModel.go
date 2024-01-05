package models

type NutritionalValue struct {
	ID		  		uint 	`gorm:"primarykey"`
	OwnerTitle		string 	`json:"owner_title"`
	OwnerID   		uint
  	OwnerType 		string
	Edited 			bool   	`json:"edited" gorm:"ignore"`

	Kcal			int    	`json:"kcal"`
	Kj				int	   	`json:"kj"`

	Fat				int    	`json:"fat"`
	SaturatedFat	int	   	`json:"saturated_fat"`
	Carbohydrate	int	   	`json:"carbohydrate"`
	Sugar			int    	`json:"sugar"`
	Protein			int    	`json:"protein"`
	Salt			int	   	`json:"salt"`
}