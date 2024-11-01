package handlers

import (
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	data := helpers.Envelope{
		"status": "available",
	}
	err := helpers.WriteJSON(w, data, http.StatusOK, nil)
	if err != nil {
		errors.ServerErrorResponse(w, err)
	}
}
