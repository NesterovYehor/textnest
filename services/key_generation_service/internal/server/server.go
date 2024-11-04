package server

import (
	"net/http"
	"time"

	"github.com/NesterovYehor/TextNest/tree/main/services/key_generation_service/internal/cache"
	"github.com/NesterovYehor/TextNest/tree/main/services/key_generation_service/internal/config"
	"github.com/NesterovYehor/TextNest/tree/main/services/key_generation_service/internal/handlers"
)

type Server struct {
	Config *config.Config
}

func (server *Server) Start() error {
	rdb, err := cache.StartRedis(server.Config.Addr)
	if err != nil {
		return err
	}
	router := http.NewServeMux()

	router.HandleFunc("GET  v1/get-key", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetKey(w, r, rdb)
	})

    router.HandleFunc("POST v1/transfer-key", func(w http.ResponseWriter, r *http.Request) {

    })
	srv := http.Server{
		Addr:         server.Config.Addr,
		Handler:      nil,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	return srv.ListenAndServe()
}
