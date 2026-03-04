package main

import (
	"log"
	"net/http"
	"time"

	"github.com/RadekKusiak71/splitwise/gateway/internal/auth"
	"github.com/RadekKusiak71/splitwise/gateway/internal/config"
	"github.com/RadekKusiak71/splitwise/gateway/internal/proxy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()

	jwtMgr := auth.NewJWTManager(cfg.JWTSecret)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Handle("/api/v1/users/*", proxy.New(cfg.IdentityServiceURL))
	r.Handle("/api/v1/auth/*", proxy.New(cfg.IdentityServiceURL))

	r.Group(func(r chi.Router) {
		r.Use(jwtMgr.JWTMiddleware)
		r.Handle("/api/v1/requests/*", proxy.New(cfg.RequestsServiceURL))
	})

	log.Printf("Gateway service is running on port :80")
	log.Printf("Routing /api/v1 to %s", cfg.IdentityServiceURL)

	server := &http.Server{
		Addr:         ":80",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
