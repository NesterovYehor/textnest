package handler

import (
	"net/http"
)

func CreateNewPost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		text string `json: "post_text`
		data string `json: "date"`
	}
}
