package handlers

import (
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg"
)

func StorePaste(w http.ResponseWriter, r *http.Request) {
	var input struct {
		hash    string `json: "hash"`
		content string `json: "content"`
	}

	err := helpers.ReadJson(w, r, &input)
}
