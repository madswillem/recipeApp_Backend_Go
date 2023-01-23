package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipeSchema = struct {
	ID 				primitive.ObjectID `bson:"_id, omitempty"`
	title			string `bson:"string field"`
	ingredients		[]IngredientsSchema `bson:"array field"`
	preparation		string `bson:"string field"`
	selected		int	   `bson:"string field"`	
	date			time.Time `bson:"date field"`
	Version     	int    `bson:"int field"`
}

type IngredientsSchema = struct {
	id				int
	ingredient		string
	amount			string
}

