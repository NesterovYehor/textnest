package errors

import (
	"fmt"
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/helpers"
)

func errorResponse(w http.ResponseWriter, status int, message any) {
	env := helpers.Envelope{
		"error": message,
	}

	err := helpers.WriteJSON(w, env, status, nil)
	if err != nil {
		w.WriteHeader(500)
	}
}

func BadRequestResponse(w http.ResponseWriter, status int, err error) {
	errorResponse(w, status, err.Error())
}

func UploadContent(w http.ResponseWriter) {
	message := fmt.Sprintln("failed to upload content to S3")
	errorResponse(w, http.StatusServiceUnavailable, message)
}

func ServerErrorResponse(w http.ResponseWriter, err error) {
	message := "the server encountered a problem and could not process your request"

	errorResponse(w, http.StatusServiceUnavailable, message)
}
