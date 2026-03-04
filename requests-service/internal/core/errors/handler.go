package errors

import (
	"log"
	"net/http"

	"github.com/RadekKusiak71/splitwise-requests/internal/core/utils"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func HandleAPIError(next APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)
		if err != nil {
			apiErr, ok := err.(*APIError)
			if ok {
				utils.WriteJSON(w, apiErr.StatusCode, apiErr)
				return
			}
			log.Printf("Unexpected APIError occured whiile processing request: %s", err.Error())
			apiErr = APIErrInternalServer
			utils.WriteJSON(w, apiErr.StatusCode, apiErr)
			return
		}
	}
}
