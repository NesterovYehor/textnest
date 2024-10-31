package main

import (
	"sync"

	"github.com/NesterovYehor/TextNest/tree/main/services/storage_service/config"
	"github.com/NesterovYehor/TextNest/tree/main/services/storage_service/server"
)

func main() {
	cfg := config.InitConfig()
	var wg sync.WaitGroup
	server.StartServer(cfg, &wg)
}
