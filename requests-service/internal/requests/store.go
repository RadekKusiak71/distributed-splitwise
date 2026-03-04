package requests

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) RequestsStore {
	return &store{db: db}
}

func (s *store) Save(ctx context.Context, r *Request) error {
	query := `
		INSERT INTO requests (
			id, user_id, idempotency_key, status, input_s3_key, output_s3_key, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			output_s3_key = EXCLUDED.output_s3_key,
			updated_at = EXCLUDED.updated_at;
	`
	_, err := s.db.Exec(ctx, query,
		r.id, r.userID, r.idempotencyKey, r.status,
		r.inputS3Key, r.outputS3Key, r.createdAt, r.updatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save request: %w", err)
	}
	return nil
}

func (s *store) GetByID(ctx context.Context, userID, requestID string) (*Request, error) {
	query := `SELECT id, user_id, idempotency_key, status, input_s3_key, output_s3_key, created_at, updated_at 
	          FROM requests WHERE id = $1 AND user_id = $2`

	return s.scanRow(s.db.QueryRow(ctx, query, requestID, userID))
}

func (s *store) GetByIdempotency(ctx context.Context, userID, idempotencyKey string) (*Request, error) {
	query := `SELECT id, user_id, idempotency_key, status, input_s3_key, output_s3_key, created_at, updated_at 
	          FROM requests WHERE idempotency_key = $1 AND user_id = $2`

	return s.scanRow(s.db.QueryRow(ctx, query, idempotencyKey, userID))
}

func (s *store) GetAll(ctx context.Context, userID string) ([]*Request, error) {
	query := `SELECT id, user_id, idempotency_key, status, input_s3_key, output_s3_key, created_at, updated_at 
	          FROM requests WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Request
	for rows.Next() {
		req, err := s.scan(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, req)
	}
	return results, nil
}

func (s *store) scanRow(row pgx.Row) (*Request, error) {
	r, err := s.scan(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return r, err
}

func (s *store) scan(scanner pgx.Row) (*Request, error) {
	r := &Request{}
	err := scanner.Scan(
		&r.id,
		&r.userID,
		&r.idempotencyKey,
		&r.status,
		&r.inputS3Key,
		&r.outputS3Key,
		&r.createdAt,
		&r.updatedAt,
	)
	if err != nil {
		return nil, err
	}
	return r, nil
}
