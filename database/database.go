package db

import (
	"fmt"
	"context"
	"time"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const databaseURI = "mongodb://localhost:27017"
const delayInSec = 10 * time.Second

type UserToDB struct {
	Email string
	PasswordHash string
}

func getDBClient() (client *mongo.Client) {
	client, err := mongo.NewClient(options.Client().ApplyURI(databaseURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), delayInSec)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")
	fmt.Println("Available dbs:", databases)

	return client
}

func SaveUserToDB(user UserToDB) (*mongo.InsertOneResult){
	var client = getDBClient()
	ctx, _ := context.WithTimeout(context.Background(), delayInSec)
	defer client.Disconnect(ctx)


	todoDatabase := client.Database("todo")
	usersCollection := todoDatabase.Collection("users")

	userSaveResult, err := usersCollection.InsertOne(ctx, bson.D{
		{Key: "email", Value: user.Email},
		{Key: "passwordHash", Value: user.PasswordHash},
	})

	if err != nil {
		fmt.Println("Failed to save new user to DB")
		log.Fatal(err)
	}


	// fmt.Println("SAVED USER RESULT", reflect.TypeOf(userSaveResult))
	return userSaveResult

}