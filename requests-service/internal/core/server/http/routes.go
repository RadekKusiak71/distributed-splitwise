package http

import (
	"log"
	"time"

	"github.com/RadekKusiak71/splitwise-requests/internal/auth"
	"github.com/RadekKusiak71/splitwise-requests/internal/core/errors"
	"github.com/RadekKusiak71/splitwise-requests/internal/core/queue"
	"github.com/RadekKusiak71/splitwise-requests/internal/core/server/middlewares"
	"github.com/RadekKusiak71/splitwise-requests/internal/core/storage"
	"github.com/RadekKusiak71/splitwise-requests/internal/requests"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (as *apiserver) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	v1 := chi.NewRouter()
	r.Mount("/api/v1", v1)

	s3Client, err := storage.NewS3Client(as.cfg.AWS.Region)
	if err != nil {
		log.Fatalf("Couldn't connect with s3 client: %s", err.Error())
	}
	uploader := storage.NewS3Uploader(s3Client, &as.cfg.AWS)

	sqsClient, err := queue.NewSQSClient(as.cfg.AWS.Region)
	if err != nil {
		log.Fatalf("Couldn't create SQS client: %s", err.Error())
	}
	publisher := queue.NewSQSPublisher(sqsClient, as.cfg.AWS.SQSQueueURL)

	requestsStore := requests.NewStore(as.db)
	requestsService := requests.NewService(requestsStore, uploader, publisher)
	requestsHandler := requests.NewHandler(requestsService)

	v1.Route("/requests", func(r chi.Router) {
		r.Use(auth.IsAuthenticated)

		// GET /api/v1/requests
		r.Get("/", errors.HandleAPIError(requestsHandler.HandleGetAllRequests))

		r.Route("/{requestID}", func(r chi.Router) {
			r.Use(requests.RequestContext)
			r.Get("/", errors.HandleAPIError(requestsHandler.HandleGetRequestByID))
		})

		r.Group(func(r chi.Router) {
			r.Use(middlewares.IdempotencyKeyRequired)
			r.Post("/", errors.HandleAPIError(requestsHandler.HandleCreateRequest))
		})
	})

	return r
}
