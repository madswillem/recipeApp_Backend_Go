package models

type UserSchema struct {
	ID           string         `json:"_id" gorm:"primarykey"`
	IP           []string       `json:"ip"`
	RecipeRating []RatingStruct `json:"rating" gorm:"polymorphic:Owner"`
}
