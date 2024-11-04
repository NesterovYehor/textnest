package handlers

import (
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/pkg/validator"
	"github.com/NesterovYehor/TextNest/tree/main/services/key_generation_service/internal/keymanager"
	"github.com/redis/go-redis/v9"
)

func TransferKey(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {
	key := r.PathValue("key")

	if key == "" {
		errors.IncorrectUrlParams(w, "key")
	}
	v := validator.New()

	if keymanager.IsKeyValid(v, key); !v.Valid() {
		errors.FailedValidationResponse(w, r, v.Errors)
	}

	err := keymanager.ReallocateKey(key, rdb)
	if err != nil {
		errors.ServerErrorResponse(w, err)
	}

	response := helpers.Envelope{
		"Status": "Succesful",
	}

	err = helpers.WriteJSON(w, response, http.StatusOK, nil)
	if err != nil {
		errors.ServerErrorResponse(w, err)
	}
}
