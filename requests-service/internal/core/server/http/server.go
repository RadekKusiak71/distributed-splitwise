package http

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RadekKusiak71/splitwise-requests/internal/core/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type apiserver struct {
	cfg *config.Config
	db  *pgxpool.Pool
}

func NewAPIServer(cfg *config.Config, db *pgxpool.Pool) *apiserver {
	return &apiserver{
		cfg: cfg,
		db:  db,
	}
}

func (as *apiserver) Run() error {
	h := &http.Server{
		Addr:         fmt.Sprintf(":%d", as.cfg.API.Port),
		Handler:      as.SetupRoutes(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Starting HTTP server on port :%d", as.cfg.API.Port)
	return h.ListenAndServe()
}
