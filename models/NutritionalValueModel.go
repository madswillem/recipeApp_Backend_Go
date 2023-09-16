package models

type NutritionalValue struct {
	ID		  		string `gorm:"primarykey"`
	OwnerTitle		string `json:"owner_title"`
	OwnerID   		string
  	OwnerType 		string

	Kcal			int    `json:"kcal"`
	Kj				int	   `json:"kj"`

	Fat				int    `json:"fat"`
	SaturatedFat	int	   `json:"saturated_fat"`
	Carbohydrate	int	   `json:"carbohydrate"`
	Sugar			int    `json:"sugar"`
	Protein			int    `json:"protein"`
	Salt			int	   `json:"salt"`
}