package middleware

import (
	"context"
	"fmt"

	"github.com/sony/gobreaker"
)

type CircuitBreakerMiddleware struct {
	breaker *gobreaker.CircuitBreaker
}

func NewCircuitBreakerMiddleware(settings gobreaker.Settings) *CircuitBreakerMiddleware {
	return &CircuitBreakerMiddleware{
		breaker: gobreaker.NewCircuitBreaker(settings),
	}
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
