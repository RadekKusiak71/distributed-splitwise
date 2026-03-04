package main

import (
	"context"
	"log"

	"github.com/RadekKusiak71/splitwise-requests/internal/core/config"
	"github.com/RadekKusiak71/splitwise-requests/internal/core/postgres"
	"github.com/RadekKusiak71/splitwise-requests/internal/core/server/http"
)

func main() {
	cfg := config.Load()
	db, err := postgres.New(context.Background(), &postgres.PGConfig{
		HOST:     cfg.DB.Host,
		PORT:     cfg.DB.Port,
		USER:     cfg.DB.User,
		PASSWORD: cfg.DB.Password,
		NAME:     cfg.DB.Name,
	})
	if err != nil {
		log.Fatalf("Couldn't connect with database: %s", err.Error())
	}

	as := http.NewAPIServer(cfg, db)
	if err := as.Run(); err != nil {
		log.Fatalf("Couldn't start HTTP server: %s", err.Error())
	}
}
