package main

import (
	"sync"

	"github.com/NesterovYehor/TextNest/services/api_service/config"
	"github.com/NesterovYehor/TextNest/services/api_service/server"
)

func main() {
	cfg := config.InitConfig()
	var wg sync.WaitGroup
	server.RunServer(cfg, &wg)
}
