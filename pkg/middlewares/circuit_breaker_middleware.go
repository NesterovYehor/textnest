package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

type CircuitBreakerConfig struct {
	MaxRequests uint32        `mapstructure:"max_requests"`
	Interval    time.Duration `mapstructure:"interval"`
	Timeout     time.Duration `mapstructure:"timeout"`
}

type CircuitBreakerMiddleware struct {
	breaker *gobreaker.CircuitBreaker
}

func NewCircuitBreakerMiddleware(cfg CircuitBreakerConfig, name string) *CircuitBreakerMiddleware {
	return &CircuitBreakerMiddleware{
		breaker: gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        name,
			MaxRequests: cfg.MaxRequests,
			Timeout:     cfg.Timeout,
			Interval:    cfg.Interval,
		}),
	}
}

func (cb *CircuitBreakerMiddleware) UpdateConfig(cfg CircuitBreakerConfig) {
	cb.breaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        cb.breaker.Name(),
		MaxRequests: cfg.MaxRequests,
		Timeout:     cfg.Timeout,
		Interval:    cfg.Interval,
	})
}

// Middleware function to wrap calls with Circuit Breaker
func (cb *CircuitBreakerMiddleware) Execute(ctx context.Context, operation func(context.Context) (any, error)) (any, error) {
	result, err := cb.breaker.Execute(func() (any, error) {
		return operation(ctx)
	})
	if err != nil {
		return nil, fmt.Errorf("circuit breaker triggered: %w", err)
	}

	return result, nil
}
