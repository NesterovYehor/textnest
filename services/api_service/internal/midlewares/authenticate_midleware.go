package midlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/app"
)

// Define a custom context key type
type contextKey string

const userIDKey contextKey = "user_id"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the AppContext instance
		appCtx, err := app.GetInstance()
		if err != nil {
			errors.ServerErrorResponse(w, err)
			return
		}

		// Get Authorization Header
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			ctx := context.WithValue(r.Context(), userIDKey, "")
			updatedRequest := r.WithContext(ctx)

			// Proceed to next handler with updated request
			next.ServeHTTP(w, updatedRequest)
		}

		// Check if it's a Bearer token
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			errors.BadRequestResponse(w, http.StatusUnauthorized, fmt.Errorf("Invalid Authentication Token Response"))
			return
		}

		// Extract token
		token := headerParts[1]

		// Authorize user
		userId, err := appCtx.AuthClient.AuthorizeUser(token)
		if err != nil {
			errors.ServerErrorResponse(w, err)
			return
		}

		// Store user ID in request context
		ctx := context.WithValue(r.Context(), userIDKey, userId)
		updatedRequest := r.WithContext(ctx)

		// Proceed to next handler with updated request
		next.ServeHTTP(w, updatedRequest)
	})
}
