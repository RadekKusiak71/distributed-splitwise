package requests

import (
	"net/http"

	"github.com/RadekKusiak71/splitwise-requests/internal/auth"
	"github.com/RadekKusiak71/splitwise-requests/internal/core/errors"
	"github.com/RadekKusiak71/splitwise-requests/internal/core/server/middlewares"
	"github.com/RadekKusiak71/splitwise-requests/internal/core/utils"
)

type handler struct {
	requestService RequestsService
}

func NewHandler(requestService RequestsService) *handler {
	return &handler{requestService: requestService}
}

func (h *handler) HandleGetAllRequests(w http.ResponseWriter, r *http.Request) error { return nil }
func (h *handler) HandleGetRequestByID(w http.ResponseWriter, r *http.Request) error { return nil }
func (h *handler) HandleCreateRequest(w http.ResponseWriter, r *http.Request) error {
	r.ParseMultipartForm(20 << 20) // 20 MB file
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		return err
	}

	idempotencyKey, err := middlewares.GetIdempotencyIDFromRequest(r)
	if err != nil {
		return err
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		return errors.APIErrBadRequest
	}

	req, err := h.requestService.CreateRequest(r.Context(), userID, idempotencyKey, file, handler)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, req)
}
