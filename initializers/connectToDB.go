package initializers

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DB *mongo.Client

func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB_CONNECTION")

	// Create a new client and connect to the server
	DB, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(dsn))
	if err != nil {
		panic(err)
	}
	// Ping the primary
	if err := DB.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
}