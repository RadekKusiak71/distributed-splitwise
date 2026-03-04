package requests

import (
	"fmt"
	"net/http"

	"github.com/RadekKusiak71/splitwise-requests/internal/core/errors"
)

var (
	ErrRequestNotFound = fmt.Errorf("error not found")
)

var (
	APIErrInvalidFileExstension       = errors.NewAPIError("Accepted file extensions are: .csv", http.StatusUnprocessableEntity)
	APIErrRequestIsScheduledToProcess = errors.NewAPIError("Request is already scheduled to process", http.StatusConflict)
)
