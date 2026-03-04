package http

import (
	"log"
	"time"

	"github.com/RadekKusiak71/splitwise-requests/internal/auth"
	"github.com/RadekKusiak71/splitwise-requests/internal/core/errors"
	"github.com/RadekKusiak71/splitwise-requests/internal/core/server/middlewares"
	"github.com/RadekKusiak71/splitwise-requests/internal/core/storage"
	"github.com/RadekKusiak71/splitwise-requests/internal/requests"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (as *apiserver) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	v1 := chi.NewRouter()

	s3Client, err := storage.NewS3Client(as.cfg.AWS.Region)
	if err != nil {
		log.Fatalf("Couldn't connect with s3 client: %s", err.Error())
	}
	uploader := storage.NewS3Uploader(s3Client, &as.cfg.AWS)

	requestsStore := requests.NewStore(as.db)
	requestsService := requests.NewService(requestsStore, uploader)
	requestsHandler := requests.NewHandler(requestsService)

	r.Mount("/api/v1", v1)

	v1.Route("/requests", func(r chi.Router) {
		r.Use(auth.IsAuthenticated)

		r.Get("/", errors.HandleAPIError(requestsHandler.HandleGetAllRequests))

		r.Group(func(r chi.Router) {
			r.Use(middlewares.IdempotencyKeyRequired)
			r.Post("/", errors.HandleAPIError(requestsHandler.HandleCreateRequest))
		})

		r.Route("/{requestID}", func(r chi.Router) {
			r.Use(requests.RequestContext)
		})

	})

	return r
}
