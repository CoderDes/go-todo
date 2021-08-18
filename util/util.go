package util

import (
	"encoding/json"
	"log"

	errConst "go-todo/constants/errors"
)

func ErrorHandler(options errConst.ErrorHandlerOptions) {
	respBody, err := json.Marshal(options.Payload)
	if err != nil {
		log.Println("Failed encoding data into JSON format")
		log.Fatal(err)
	}

	options.RespWr.Write(respBody)
}
