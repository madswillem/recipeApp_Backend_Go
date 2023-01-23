package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipeSchema = struct {
	ID 				primitive.ObjectID `bson:"_id, omitempty"`
	Title			string `bson:"title, omitempty"`
	Ingredients		[]IngredientsSchema `bson:"ingredients, omitempty"`
	Preparation		string `bson:"preparation, omitempty"`
	Selected		int	   `bson:"selected, omitempty"`	
	Date			time.Time `bson:"date, omitempty"`
	Version     	int    `bson:"__v, omitempty"`
}

type IngredientsSchema = struct {
	Id				string `bson:"id"`
	Ingredient		string `bson:"ingredient"`
	Amount			string `bson:"amount"`
}

