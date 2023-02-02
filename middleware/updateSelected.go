package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rezeptapp.ml/goApp/initializers"
)

func UpdateSelected(id string, change int, c *gin.Context)(*mongo.UpdateResult) {
	coll := initializers.DB.Database("test").Collection("recepies")

	result := GetDataByID(id, c)

	result.Selected += change

	filter := bson.D{{"_id", result.ID}}
	update := bson.D{{"$set", bson.D{{"selected", result.Selected}}}}
	res, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	return res
}