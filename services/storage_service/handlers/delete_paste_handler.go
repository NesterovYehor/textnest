package handlers

import (
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/tree/main/services/storage_service/models"
)

func DeletePaste(w http.ResponseWriter, r *http.Request, storage models.Storage) {
	key := r.PathValue("key")

	if key == "" {
		errors.IncorrectUrlParams(w, "key")
		return
	}

	err := storage.DeletePaste(key)
	if err != nil {
		errors.ServerErrorResponse(w, err)
		return
	}

	response := helpers.Envelope{
		"status": "successfuly deleted content",
	}

	err = helpers.WriteJSON(w, response, http.StatusOK, nil)
	if err != nil {
		errors.ServerErrorResponse(w, err)
		return
	}
}
