package constants

import "net/http"

type ErrorResponseMessage struct {
	StatusCode int
	ErrorMsg   string
}

type ErrorHandlerOptions struct {
	RespWr  http.ResponseWriter
	Payload ErrorResponseMessage
}
