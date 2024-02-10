package endpoints

import (
	"context"
	"emailn/internal/infrastructure/credential"
	"net/http"

	"github.com/go-chi/render"
)

type ValidateTokenFunc func(token string, ctx context.Context) (string, error)

var ValidateToken ValidateTokenFunc = credential.ValidateToken

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			render.Status(r, 401)
			render.JSON(w, r, map[string]string{"error": "request does not contain an authorization header"})
			return
		}

		email, err := ValidateToken(tokenString, r.Context())
		if err != nil {
			render.Status(r, 401)
			render.JSON(w, r, map[string]string{"error": "invalid token"})
			return
		}

		ctx := context.WithValue(r.Context(), "email", email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
