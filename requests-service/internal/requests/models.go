package requests

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type RequestStatus string

const (
	Pending    RequestStatus = "pending"
	Processing RequestStatus = "processing"
	Completed  RequestStatus = "completed"
	Failed     RequestStatus = "failed"
)

type Request struct {
	id             string
	userID         string
	idempotencyKey string
	status         RequestStatus
	inputS3Key     string
	outputS3Key    sql.NullString
	createdAt      time.Time
	updatedAt      time.Time
}

func NewRequest(userID, idempotencyKey, inputS3Key string) *Request {
	now := time.Now()
	return &Request{
		id:             uuid.NewString(),
		userID:         userID,
		idempotencyKey: idempotencyKey,
		status:         Pending,
		inputS3Key:     inputS3Key,
		createdAt:      now,
		updatedAt:      now,
	}
}

func (r *Request) ID() string                  { return r.id }
func (r *Request) UserID() string              { return r.userID }
func (r *Request) IdempotencyKey() string      { return r.idempotencyKey }
func (r *Request) Status() RequestStatus       { return r.status }
func (r *Request) InputS3Key() string          { return r.inputS3Key }
func (r *Request) OutputS3Key() sql.NullString { return r.outputS3Key }
func (r *Request) CreatedAt() time.Time        { return r.createdAt }
func (r *Request) UpdatedAt() time.Time        { return r.updatedAt }
