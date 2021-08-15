package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const databaseURI = "mongodb://localhost:27017"
const delayInSec = 10 * time.Second
const serverPort = ":8080"

type User struct {
	Email    string `json: "email"`
	Password string `json: "password"`
}

func loginHandler(rewWr http.ResponseWriter, req *http.Request) {
	fmt.Println("LOGIN REQUEST IS ", req)
}

func registerHandler(resWr http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(resWr, "Method is not supported", http.StatusNotFound)
		return
	}
	// TODO: check Content-type header for application/json

	user := User{}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: hash the password and write user to DB

	fmt.Println("USER", user)
}

func getDBConnection() (client *mongo.Client) {
	client, err := mongo.NewClient(options.Client().ApplyURI(databaseURI))
	if err != nil {
		log.Fatal(err)
	}

	// what for this
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

	fmt.Println("Connected to DB")
	fmt.Println("dbs:", databases)

	return client
}

func main() {
	var dbClient = getDBConnection()
	ctx, _ := context.WithTimeout(context.Background(), delayInSec)
	defer dbClient.Disconnect(ctx)

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)

	if err := http.ListenAndServe(serverPort, nil); err != nil {
		log.Fatal(err)
	}
}
