package models

type DietSchema struct {
	ID         uint   `gorm:"primarykey"`
	OwnerTitle string `json:"owner_title"`
	OwnerID    string
	OwnerType  string

	Vegetarien    bool `json:"vegetarien"`
	Vegan         bool `json:"vegan"`
	LowCal        bool `json:"lowcal"`
	LowCarb       bool `json:"lowcarb"`
	Keto          bool `json:"keto"`
	Paleo         bool `json:"paleo"`
	LowFat        bool `json:"lowfat"`
	FoodCombining bool `json:"food_combining"`
	WholeFood     bool `json:"whole_food"`
}
