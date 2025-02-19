package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/services/api_service/config"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func RateLimit(cfg *config.Config, next http.Handler) http.Handler {
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Start background cleanup goroutine
	if cfg.Limiter.Enabled {
		go cleanupOldEntries(&mu, clients, cfg.Limiter.CleanupInterval)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !cfg.Limiter.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		ip, err := getClientIP(r)
		if err != nil {
			errors.BadRequestResponse(w, http.StatusInternalServerError, fmt.Errorf("Could not validate client address"))
			return
		}

		mu.Lock()
		defer mu.Unlock()

		// Create new client if not exists
		if _, exists := clients[ip]; !exists {
			clients[ip] = &client{
				limiter: rate.NewLimiter(
					rate.Limit(cfg.Limiter.Rps),
					cfg.Limiter.Burst,
				),
			}
		}

		// Update last seen time
		clients[ip].lastSeen = time.Now()

		// Enforce rate limit
		if !clients[ip].limiter.Allow() {
			errors.RateLimitExceededResponse(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getClientIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	return ip, nil
}

func cleanupOldEntries(mu *sync.Mutex, clients map[string]*client, interval time.Duration) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		mu.Lock()
		for ip, c := range clients {
			if time.Since(c.lastSeen) > interval {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}
