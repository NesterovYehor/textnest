package handlers

import (
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/validator"
	"github.com/NesterovYehor/TextNest/tree/main/services/key_generation_service/internal/keymanager"
	"github.com/redis/go-redis/v9"
)

func transferKey(w http.ResponseWriter, r http.Request, rdb *redis.Client){
    key := r.PathValue("key")
    
    if key == ""{
        errors.IncorrectUrlParams(w, "key")
    }
    v := validator.New()

    if keymanager.IsKeyValid(v, key); !v.Valid(){
        

    }
    
}
