package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/models"
)

func GetDataByID(id string, c *gin.Context)(models.RecipeSchema) {
	coll := initializers.DB.Database("test").Collection("recepies")

	// Declare Context type object for managing multiple API requests
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// convert id string to ObjectId
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	}

	// find
	result := models.RecipeSchema{}
	err = coll.FindOne(ctx, bson.M{"_id": objectId}).Decode(&result)

	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	}

	return result
}