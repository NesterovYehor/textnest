package errors

import (
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/helpers"
)

type envelope map[string]any

func errorResponse(w http.ResponseWriter, status int, message any) {
	env := envelope{
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
