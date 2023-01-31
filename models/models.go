package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipeSchema = struct {
	ID 				primitive.ObjectID	`bson:"_id, omitempty" json:"_id"`
	Title			string				`bson:"title, omitempty" json:"title"` 
	Ingredients		[]IngredientsSchema	`bson:"ingredients, omitempty" json:"ingredients"`
	Preparation		string				`bson:"preparation, omitempty" json:"preparation"`
	Selected		int					`bson:"selected, omitempty" json:"selected"`	
	Date			time.Time			`bson:"date, omitempty" json:"date"`
	Version     	int					`bson:"__v, omitempty" json:"__v"`
}

type IngredientsSchema = struct {
	Id				string				`bson:"id" json:"id"` 
	Ingredient		string				`bson:"ingredient" json:"ingredient"` 
	Amount			string				`bson:"amount" json:"amount"` 
}