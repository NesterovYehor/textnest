package handlers

import (
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/tree/main/services/key_generation_service/internal/keymanager"

	"github.com/redis/go-redis/v9"
)

func GetKey(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {
	key, err := keymanager.GetKey(rdb)
	if err != nil {
		errors.ServerErrorResponse(w, err)
		return
	}

	response := helpers.Envelope{
		"key": key,
	}

	err = helpers.WriteJSON(w, response, http.StatusOK, nil)
	if err != nil {
		errors.ServerErrorResponse(w, err)
		return
	}
}
