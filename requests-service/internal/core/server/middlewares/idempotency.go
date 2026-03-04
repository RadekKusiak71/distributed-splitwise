package middlewares

import (
	"context"
	"net/http"

	"github.com/RadekKusiak71/splitwise-requests/internal/core/errors"
	"github.com/RadekKusiak71/splitwise-requests/internal/core/utils"
	"github.com/google/uuid"
)

type contextKey string

const RequestIdempotencyID contextKey = "idempodency_key"

func IdempotencyKeyRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idempotency := r.Header.Get("X-Idempotency-Key")
		if idempotency == "" {
			err := errors.APIErrMissingIdempotencyKey
			utils.WriteJSON(w, err.StatusCode, err)
			return
		}

		idempotencyUUID, err := uuid.Parse(idempotency)
		if err != nil {
			err := errors.APIErrMissingIdempotencyKey
			utils.WriteJSON(w, err.StatusCode, err)
			return
		}

		ctx := context.WithValue(r.Context(), RequestIdempotencyID, idempotencyUUID.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetIdempotencyIDFromRequest(r *http.Request) (string, error) {
	idempotencyID, ok := r.Context().Value(RequestIdempotencyID).(string)
	if !ok {
		return "", errors.APIErrInternalServer
	}
	return idempotencyID, nil
}
