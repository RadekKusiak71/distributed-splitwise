package requests

import (
	"context"
	"net/http"

	"github.com/RadekKusiak71/splitwise-requests/internal/core/errors"
	"github.com/RadekKusiak71/splitwise-requests/internal/core/utils"
	"github.com/go-chi/chi"
)

type contextKey string

const RequestID contextKey = "request_id"

func RequestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "requestID")

		if id == "" {
			err := errors.APIErrBadRequest
			utils.WriteJSON(w, err.StatusCode, err)
			return
		}

		ctx := context.WithValue(r.Context(), RequestID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestIDFromRequest(r *http.Request) (string, error) {
	requestID, ok := r.Context().Value(RequestID).(string)
	if !ok {
		return "", errors.APIErrInternalServer
	}
	return requestID, nil
}
