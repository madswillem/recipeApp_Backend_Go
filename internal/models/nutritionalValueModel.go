package models

type NutritionalValue struct {
	ID         uint   `gorm:"primarykey"`
	OwnerTitle string `json:"owner_title"`
	OwnerID    uint
	OwnerType  string
	Edited     bool `json:"edited" gorm:"ignore"`

	Kcal float64 `json:"kcal"`
	Kj   float64 `json:"kj"`

	Fat          float64 `json:"fat"`
	SaturatedFat float64 `json:"saturated_fat"`
	Carbohydrate float64 `json:"carbohydrate"`
	Sugar        float64 `json:"sugar"`
	Protein      float64 `json:"protein"`
	Salt         float64 `json:"salt"`
}
