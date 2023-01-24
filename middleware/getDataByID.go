package middleware

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/models"
)

func GetDataByID(id string)(models.RecipeSchema) {
	coll := initializers.DB.Database("test").Collection("recepies")

	// Declare Context type object for managing multiple API requests
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// convert id string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	// find
	result := models.RecipeSchema{}
	err = coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(&result)

	if err != nil {
		log.Fatal("FindOne() ObjectIDFromHex ERROR:", err)
	}

	return result
}