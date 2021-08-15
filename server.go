package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const databaseURI = "mongodb://localhost:27017"
const delayInSec  = 10 * time.Second
const serverPort = ":8080"


func loginHandler(rewWr http.ResponseWriter, req *http.Request) {
	fmt.Println("LOGIN monREQUEST IS ", req)
}

func registerHandler(resWr http.ResponseWriter, req *http.Request) {
	fmt.Println("REGISTER REQUEST IS ", req)
	if req.Method != "POST" {
		http.Error(resWr, "Method is not supported", http.StatusNotFound)
		return
	}
}	

func main() {
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

	fmt.Println(databases)

	defer client.Disconnect(ctx)


	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)

	if err := http.ListenAndServe(serverPort, nil); err != nil {
		log.Fatal(err)
	}
}