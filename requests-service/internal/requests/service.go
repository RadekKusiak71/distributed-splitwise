package requests

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"slices"

	"github.com/google/uuid"
)

var ValidExtensions = []string{".csv"}

type service struct {
	requestStore RequestsStore
	uploader     FileUploader
	publisher    MessagePublisher
}

func NewService(requestStore RequestsStore, uploader FileUploader, publisher MessagePublisher) RequestsService {
	return &service{requestStore: requestStore, uploader: uploader, publisher: publisher}
}

func (s *service) GetAllRequests(ctx context.Context, userID string) ([]*RequestResponse, error) {
	reqs, err := s.requestStore.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	var res []*RequestResponse
	for _, r := range reqs {
		res = append(res, s.mapToResponse(r))
	}
	return res, nil
}

func (s *service) GetRequestByID(ctx context.Context, userID, requestID string) (*RequestResponse, error) {
	req, err := s.requestStore.GetByID(ctx, userID, requestID)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.New("request not found")
	}
	return s.mapToResponse(req), nil
}

func (s *service) CreateRequest(ctx context.Context, userID, idempotencyKey string, file multipart.File, fileHeaders *multipart.FileHeader) (*CreateRequestResponse, error) {
	existing, _ := s.requestStore.GetByIdempotency(ctx, userID, idempotencyKey)
	if existing != nil {
		return nil, errors.New("request already exists (idempotency conflict)")
	}

	fileExt := filepath.Ext(fileHeaders.Filename)
	if !slices.Contains(ValidExtensions, fileExt) {
		return nil, errors.New("invalid file extension")
	}

	fileKey := fmt.Sprintf("uploads/%s/%s%s", userID, uuid.New().String(), fileExt)
	if err := s.uploader.Upload(ctx, fileKey, file); err != nil {
		return nil, err
	}

	req := NewRequest(userID, idempotencyKey, fileKey)
	if err := s.requestStore.Save(ctx, req); err != nil {
		return nil, err
	}

	msg := map[string]string{
		"request_id":  req.id,
		"s3_file_key": req.inputS3Key,
	}
	if err := s.publisher.Publish(ctx, msg); err != nil {
		return nil, fmt.Errorf("publish sqs message: %w", err)
	}

	return &CreateRequestResponse{
		ID:             req.id,
		IdempotencyKey: req.idempotencyKey,
		Status:         req.status,
		FileLink:       s.uploader.GenerateObjectURL(req.inputS3Key),
		CreatedAt:      req.createdAt,
	}, nil
}

func (s *service) mapToResponse(r *Request) *RequestResponse {
	var outLinkPtr *string

	if r.outputS3Key.Valid && r.outputS3Key.String != "" {
		link := s.uploader.GenerateObjectURL(r.outputS3Key.String)
		outLinkPtr = &link
	}

	return &RequestResponse{
		ID:             r.id,
		Status:         r.status,
		InputFileLink:  s.uploader.GenerateObjectURL(r.inputS3Key),
		OutputFileLink: outLinkPtr,
		CreatedAt:      r.createdAt,
		UpdatedAt:      r.updatedAt,
	}
}
