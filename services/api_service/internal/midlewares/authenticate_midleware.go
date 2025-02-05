package midlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/app"
)



func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appCtx, err := app.GetInstance()
		if err != nil {
			errors.ServerErrorResponse(w, err)
			return
		}

		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
            ctx := context.WithValue(r.Context(), "user_id", "")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			errors.BadRequestResponse(w, http.StatusUnauthorized, fmt.Errorf("Invalid Authentication Token Response"))
			return
		}

		token := headerParts[1]

		userId, err := appCtx.AuthClient.AuthorizeUser(token)
		if err != nil {
			appCtx.Logger.PrintError(context.Background(), fmt.Errorf("Failed to authorize user: %v", err), nil)
			errors.ServerErrorResponse(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userId)
        fmt.Println(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
