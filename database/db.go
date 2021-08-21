package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	dbConst "go-todo/constants/db"
	usrConst "go-todo/constants/user"
)

func getDBClient() (client *mongo.Client) {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbConst.DatabaseURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), dbConst.DelayInSec)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func UserAlreadyExists(email string) (exists bool) {
	_, err := getUserByEmail(email)
	if err != nil {
		return false
	}

	return true
}

func getUserByEmail(email string) (userFromDB usrConst.UserToDB, err error) {
	var client = getDBClient()
	ctx, _ := context.WithTimeout(context.Background(), dbConst.DelayInSec)
	defer client.Disconnect(ctx)

	todoDatabase := client.Database("todo")
	usersCollection := todoDatabase.Collection("users")

	var user usrConst.UserToDB
	
	filter := bson.M{"email": email}
	cursor := usersCollection.FindOne(ctx, filter)
	err = cursor.Decode(&user)
	
	return user, err
}

func SaveUserToDB(user usrConst.UserToDB) *mongo.InsertOneResult {
	var client = getDBClient()
	ctx, _ := context.WithTimeout(context.Background(), dbConst.DelayInSec)
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

	return userSaveResult
}
