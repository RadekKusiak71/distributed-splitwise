package auth

import (
	"context"
	"net/http"

	"github.com/RadekKusiak71/splitwise-requests/internal/core/errors"
	"github.com/RadekKusiak71/splitwise-requests/internal/core/utils"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("X-User-ID")
		apiErr := errors.APIErrUnathorized
		if userId == "" {
			utils.WriteJSON(w, apiErr.StatusCode, apiErr)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromRequest(r *http.Request) (string, error) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		return "", errors.APIErrInternalServer
	}
	return userID, nil
}
