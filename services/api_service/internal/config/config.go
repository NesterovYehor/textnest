package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/NesterovYehor/TextNest/pkg/http"
)

var (
	version   string
	buildTime string
)

type Config struct {
	Env  string
	Http *httpserver.Config
	Grpc struct {
		Addr string
	}
}

func InitConfig() *Config {
	cfg := &Config{}
    port := ":8989"
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.Grpc.Addr, "localhost:5555", "localhost:5555", "localhost:5555")

	cfg.Http = httpserver.NewConfig(port)

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Service Version:\t%s\n", version)
		fmt.Printf("Build Time:\t%s\n", buildTime)
		os.Exit(0)
	}

	return cfg
}
