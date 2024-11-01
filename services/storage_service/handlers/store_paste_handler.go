package handlers

import (
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/tree/main/services/storage_service/models"
)

func StorePaste(w http.ResponseWriter, r *http.Request, storage models.Storage) {
	var input struct {
		Key     string `json:"key"`     // Corrected tag
		Content string `json:"content"` // Corrected tag
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		errors.BadRequestResponse(w, http.StatusBadRequest, err)
		return
	}

	data := models.PasteData{
		Key:     input.Key,
		Content: input.Content,
	}

	contentPath, err := storage.UploadPaste(&data)
	if err != nil {
		errors.UploadContent(w, err) // Consider adding status code
		return
	}

	responseData := helpers.Envelope{
		"content_path": contentPath,
	}

	err = helpers.WriteJSON(w, responseData, http.StatusOK, nil)
	if err != nil {
		errors.ServerErrorResponse(w, err)
	}
}
