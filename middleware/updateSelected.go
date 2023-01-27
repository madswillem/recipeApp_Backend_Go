package middleware

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rezeptapp.ml/goApp/initializers"
)

func UpdateSelected(id string, change int)(*mongo.UpdateResult) {
	coll := initializers.DB.Database("test").Collection("recepies")

	result := GetDataByID(id)	

	result.Selected += change

	filter := bson.D{{"_id", result.ID}}
	update := bson.D{{"$set", bson.D{{"selected", result.Selected}}}}
	res, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	return res
}