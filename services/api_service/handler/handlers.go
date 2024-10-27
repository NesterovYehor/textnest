package handler

import (
	"net/http"
	"github.com/NesterovYehor/pastebin/pkg/helpers"
)

func CreateNewPost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		text string `json: "post_text`
		data string `json: "date"`
	}

	err := helper
}
