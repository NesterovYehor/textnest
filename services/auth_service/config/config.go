package config

import (
	"errors"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JwtConfig struct {
	AccessSecret     string        `mapstructure:"access_secret"`
	RefreshSecret    string        `mapstructure:"refresh_secret"`
	ActivateSecret   string        `mapstructure:"activate_secret"`
	AccessExpiry     time.Duration `mapstructure:"access_expiry"`
	RefreshExpiry    time.Duration `mapstructure:"refresh_expiry"`
	ActivateExpiry   time.Duration `mapstructure:"activate_expiry"`
	SigningMethod    jwt.SigningMethod
	SigningMethodStr string `mapstructure:"signing_method"` // Temporary field to store the raw string
}

type MailerConfig struct {
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Sender   string `mapstructure:"sender"`
	Port     int    `mapstructure:"port"`
}

type Config struct {
	DBUrl     string           `mapstructure:"db_url"`
	Grpc      *grpc.GrpcConfig `mapstructure:"grpc"`
	JwtConfig *JwtConfig       `mapstructure:"jwt"`
	Mailer    *MailerConfig    `mapstructure:"mailer"`
}

// LoadConfig reads the configuration from the config.yaml file and unmarshals it into the Config struct.
func LoadConfig(log *jsonlog.Logger) (*Config, error) {
	// Initialize Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/app") // Look in the /app directory inside the container

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshal config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Validate required fields
	if config.DBUrl == "" {
		return nil, errors.New("database URL is not provided")
	}
	if config.JwtConfig.AccessSecret == "" || config.JwtConfig.RefreshSecret == "" {
		return nil, errors.New("JWT secrets are not provided")
	}

	// Convert signing method from string
	signingMethod, err := getSigningMethod(config.JwtConfig.SigningMethodStr)
	if err != nil {
		return nil, err
	}
	config.JwtConfig.SigningMethod = signingMethod

	return &config, nil
}

func getSigningMethod(algorithm string) (jwt.SigningMethod, error) {
	if algorithm == "" {
		return nil, errors.New("JWT signing method is not set")
	}

	switch algorithm {
	case "HS256":
		return jwt.SigningMethodHS256, nil
	case "HS384":
		return jwt.SigningMethodHS384, nil
	case "HS512":
		return jwt.SigningMethodHS512, nil
	case "RS256":
		return jwt.SigningMethodRS256, nil
	case "RS384":
		return jwt.SigningMethodRS384, nil
	case "RS512":
		return jwt.SigningMethodRS512, nil
	case "ES256":
		return jwt.SigningMethodES256, nil
	default:
		return nil, errors.New("unsupported signing algorithm: " + algorithm)
	}
}
