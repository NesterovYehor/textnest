package config

import (
	"flag"
	"fmt"
	"os"
)

var (
	version   string // These variables can be set during build time
	buildTime string
)

type Config struct {
	addr string
	env  string
}

func (cfg *Config) LoadConfig() {
	flag.StringVar(&cfg.addr, "port", "8080", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Service 2 Version:\t%s\n", version)
		fmt.Printf("Build Time:\t%s\n", buildTime)
		os.Exit(0)
	}
}
