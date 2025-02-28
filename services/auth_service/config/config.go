package config

import (
	"errors"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/grpc"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	middleware "github.com/NesterovYehor/TextNest/pkg/middlewares"
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
	SigningMethodStr string        `mapstructure:"signing_method"`
	SigningMethod    jwt.SigningMethod
}

type MailerConfig struct {
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Sender   string `mapstructure:"sender"`
	Port     int    `mapstructure:"port"`
}

type DBConfig struct {
	Link            string        `mapstructure:"link"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
}

type Config struct {
	DB        *DBConfig                        `mapstructure:"db"`
	Grpc      *grpc.GrpcConfig                 `mapstructure:"grpc"`
	JwtConfig *JwtConfig                       `mapstructure:"jwt"`
	Mailer    *MailerConfig                    `mapstructure:"mailer"`
	CBConfig  *middleware.CircuitBreakerConfig `mapstructure:"circuit_breacker"`
}

func LoadConfig(log *jsonlog.Logger) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/app") // Ensure this path is correct or adjust if necessary

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshal config into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Validate required fields
	if config.DB.Link == "" {
		return nil, errors.New("no link to database provided")
	}
	if config.DB.MaxIdleConns == 0 {
		config.DB.MaxIdleConns = 25 // Set default
	}
	if config.DB.MaxOpenConns == 0 {
		config.DB.MaxOpenConns = 25 // Set default, previously incorrectly set to MaxIdleConns
	}
	if config.DB.ConnMaxLifetime == 0 {
		config.DB.ConnMaxLifetime = time.Minute * 5 // Set default
	}

	// Convert signing method string to the correct jwt.SigningMethod
	signingMethod, err := getSigningMethod(config.JwtConfig.SigningMethodStr)
	if err != nil {
		return nil, err
	}
	config.JwtConfig.SigningMethod = signingMethod

	return &config, nil
}

// Helper function to convert string to the appropriate jwt.SigningMethod
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
