package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGConfig struct {
	HOST     string
	PORT     int
	USER     string
	PASSWORD string
	NAME     string
}

func GetDBString(host, user, password, name string, port int) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, name)
}

func New(ctx context.Context, config *PGConfig) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(ctx, GetDBString(
		config.HOST, config.USER,
		config.PASSWORD, config.NAME, config.PORT,
	))
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}
