package config

type Config struct {
    Addr int
}

func InitConfig() *Config {
	cfg := &Config{} // Initialize cfg
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg.Addr = os.Getenv("PORT")
	cfg.Storage.Bucket = os.Getenv("BUCKET") // Correct key
	return cfg
}

