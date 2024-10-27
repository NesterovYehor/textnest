package main

import (
	"sync"

	"github.com/NesterovYehor/pastebin/tree/main/internal/services/api_service/config"
	"github.com/NesterovYehor/pastebin/tree/main/internal/services/api_service/server"
)

func main() {
	cfg := config.InitConfig()
	var wg sync.WaitGroup
	server.RunServer(cfg, &wg)
}
