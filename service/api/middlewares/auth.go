package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/evaevangelisti/wasatext/service/api/repositories"
	"github.com/evaevangelisti/wasatext/service/utils/errors"
	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthMiddleware(userRepository *repositories.UserRepository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errors.WriteHTTPError(w, errors.ErrUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			errors.WriteHTTPError(w, errors.ErrUnauthorized)
			return
		}

		userID := parts[1]

		uid, err := uuid.Parse(userID)
		if err != nil {
			errors.WriteHTTPError(w, errors.ErrUnauthorized)
			return
		}

		user, err := userRepository.GetUserByID(uid)
		if err != nil {
			errors.WriteHTTPError(w, err)
			return
		}

		if user == nil {
			errors.WriteHTTPError(w, errors.ErrUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}
