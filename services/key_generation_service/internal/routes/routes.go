// internal/routes/routes.go
package routes

import (
	"net/http"

	"github.com/NesterovYehor/TextNest/tree/main/services/key_generation_service/internal/handlers"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(mux *http.ServeMux, rdb *redis.Client) {
	mux.HandleFunc("/v1/get-key", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.GetKey(w, r, rdb)
	})

	mux.HandleFunc("/v1/transfer-key", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.TransferKey(w, r, rdb)
	})
}
