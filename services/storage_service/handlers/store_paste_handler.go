package handlers

import (
	"fmt"
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
)

func StorePaste(w http.ResponseWriter, r *http.Request) {
	var input struct {
		hash    string `json: "hash"`
		content string `json: "content"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
	}
}
