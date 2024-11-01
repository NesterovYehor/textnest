package handlers

import (
	"fmt"
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/tree/main/services/storage_service/models"
)

func FetchPaste(w http.ResponseWriter, r *http.Request, storage models.Storage) {
	// Corrected: Key field is now exported so JSON unmarshalling works.
	key := r.PathValue("key")

	if key == "" {
		errors.IncorrectUrlParams(w, "key")
		return
	}

	// Retrieve the paste from storage
	content, err := storage.GetPaste(key)
	if err != nil {
		fmt.Println("Error fetching paste:", err)
		errors.ServerErrorResponse(w, err)
		return
	}

	// Prepare the response data
	data := helpers.Envelope{
		"paste_content": content,
	}

	// Send JSON response with status OK (200)
	err = helpers.WriteJSON(w, data, http.StatusOK, nil) // Status code moved to correct position
	if err != nil {
		errors.ServerErrorResponse(w, err)
	}
}
