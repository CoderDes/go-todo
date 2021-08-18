package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	errConst "go-todo/constants/errors"
	usrConst "go-todo/constants/user"
	"go-todo/database"
	"go-todo/util"
)

const serverPort = ":8080"
const hashSalt = 15

type RegisterResponse struct {
	StatusCode    int
	CreatedUserId string
}

func hashPassword(password string) string {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), hashSalt)
	if err != nil {
		fmt.Println("Hashing password is failed")
		log.Fatal(err)
	}
	return string(hashBytes)
}

func loginHandler(rewWr http.ResponseWriter, req *http.Request) {
	fmt.Println("LOGIN REQUEST IS ", req)
}

func registerHandler(resWr http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		util.ErrorHandler(errConst.ErrorHandlerOptions{
			RespWr: resWr,
			Payload: errConst.ErrorResponseMessage{
				StatusCode: http.StatusBadRequest,
				ErrorMsg:   "Method is not supported",
			},
		})
		return
	}
	if req.Header.Get("Content-Type") != "application/json" {
		util.ErrorHandler(errConst.ErrorHandlerOptions{
			RespWr: resWr,
			Payload: errConst.ErrorResponseMessage{
				StatusCode: http.StatusBadRequest,
				ErrorMsg:   "Request payload must be in JSON format",
			},
		})
		return 
	}

	user := usrConst.UserFromJSON{}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	hashedPassword := hashPassword(user.Password)

	newUser := usrConst.UserToDB{
		Email:        user.Email,
		PasswordHash: hashedPassword,
	}

	// TODO: register only users with unique email
	saveResult := db.SaveUserToDB(newUser)
	userId := saveResult.InsertedID.(primitive.ObjectID).Hex()

	responseData := RegisterResponse{
		StatusCode:    http.StatusCreated,
		CreatedUserId: userId,
	}

	jsonResponseData, err := json.Marshal(responseData)
	if err != nil {
		log.Fatal(err)
	}

	resWr.Header().Set("Content-Type", "application/json")
	resWr.WriteHeader(http.StatusCreated)
	resWr.Write(jsonResponseData)
}

func main() {

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)

	// TODO: log somehow that server start listening
	if err := http.ListenAndServe(serverPort, nil); err != nil {
		log.Fatal(err)
	}
}
