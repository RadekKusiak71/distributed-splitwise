package requests

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type RequestResponse struct {
	ID             string        `json:"id"`
	Status         RequestStatus `json:"status"`
	InputFileLink  string        `json:"input_file_link"`
	OutputFileLink *string       `json:"output_file_link,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type CreateRequestResponse struct {
	ID             string        `json:"id"`
	IdempotencyKey string        `json:"idempotency_key"`
	Status         RequestStatus `json:"status"`
	FileLink       string        `json:"file_link"`
	CreatedAt      time.Time     `json:"created_at"`
}

type RequestsHandler interface {
	HandleGetAllRequests(w http.ResponseWriter, r *http.Request) error
	HandleGetRequestByID(w http.ResponseWriter, r *http.Request) error
	HandleCreateRequest(w http.ResponseWriter, r *http.Request) error
}

type RequestsService interface {
	GetAllRequests(ctx context.Context, userID string) ([]*RequestResponse, error)
	GetRequestByID(ctx context.Context, userID, requestID string) (*RequestResponse, error)
	CreateRequest(ctx context.Context, userID, idempotencyKey string, file multipart.File, fileHeaders *multipart.FileHeader) (*CreateRequestResponse, error)
}

type RequestsStore interface {
	GetAll(ctx context.Context, userID string) ([]*Request, error)
	GetByID(ctx context.Context, userID, requestID string) (*Request, error)
	GetByIdempotency(ctx context.Context, userID, idempotencyKey string) (*Request, error)
	Save(ctx context.Context, request *Request) error
}

type FileUploader interface {
	Upload(ctx context.Context, key string, file io.Reader) error
	GenerateObjectURL(key string) string
}

type MessagePublisher interface {
	Publish(ctx context.Context, message any) error
}
