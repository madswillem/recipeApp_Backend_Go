package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rezeptapp.ml/goApp/initializers"
	"rezeptapp.ml/goApp/models"
)

func GetAll(c *gin.Context) {
	var results []bson.M
	coll := initializers.DB.Database("test").Collection("recepies")
	cursor, err := coll.Find(context.TODO(), bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "Failed to read body",
			"errMessage": err.Error(),
		})
		return
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, results)

	defer cursor.Close(context.TODO())
}

func AddRecipe(c *gin.Context) {

	coll := initializers.DB.Database("test").Collection("recepies")

	type Ingredients struct {
		Id         string `json:"id"`
		Ingredient string `json:"ingredient"`
		Amount     string `json:"amount"`
	}

	var body struct {
		Title       string        `json:"title"`
		Ingredients []Ingredients `json:"ingredients"`
		Preparation string        `json:"preparation"`
		Selected    int       
		Date        time.Time
		Version     int 		  `bson:"__v"`
	}

	err := c.Bind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Failed to read body",
			"errMessage": err.Error(),
		})
		return
	}

	body.Date = time.Now()

	result, err := coll.InsertOne(context.TODO(), body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Failed to read body",
			"errMessage": err.Error(),
			"result":     result,
		})
		return
	}

	c.JSON(http.StatusCreated, body)
}

func getDataByID(id string)(models.RecipeSchema) {
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

	return  result
}

func GetById(c *gin.Context)  {
	result := getDataByID(c.Param("id"))

	c.JSON(http.StatusOK, gin.H{
		"_id": result.ID,
		"title": result.Title,
		"ingredeants": result.Ingredients,
		"preparation": result.Preparation,
		"selected": result.Selected,
		"date": result.Date,
		"__v": result.Version,
	})
}

func Select(c *gin.Context) {
	coll := initializers.DB.Database("test").Collection("recepies")

	result:= getDataByID(c.Param("id"))	

	result.Selected += 1

	filter := bson.D{{"_id", result.ID}}
	update := bson.D{{"$set", bson.D{{"selected", result.Selected}}}}
	res, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, res)
}