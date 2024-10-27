package config

import (
	"flag"
	"fmt"
	"os"
)

var (
	version   string
	buildTime string
)

type Config struct {
	Addr string
	Env  string
}

func InitConfig() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.Addr, "port", "8080", "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Service Version:\t%s\n", version)
		fmt.Printf("Build Time:\t%s\n", buildTime)
		os.Exit(0)
	}

	return cfg
}
