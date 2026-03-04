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

var (
	ValidExstenstions = []string{".csv"}
)

type service struct {
	requestStore RequestsStore
	uploader     FileUploader
}

func NewService(requestStore RequestsStore, uploader FileUploader) *service {
	return &service{requestStore: requestStore, uploader: uploader}
}

func (s *service) GetAllRequests(ctx context.Context, userID string) ([]Request, error) {
	return nil, nil
}
func (s *service) GetRequestByID(ctx context.Context, userID string, requestID string) (*Request, error) {
	return nil, nil
}
func (s *service) CreateRequest(
	ctx context.Context,
	userID, idempotencyKey string,
	file multipart.File,
	fileHeaders *multipart.FileHeader,
) (*CreateRequestResponse, error) {
	if idempotencyKey, err := s.requestStore.GetByIdempotency(ctx, userID, idempotencyKey); idempotencyKey != nil {
		if err != nil && !errors.Is(err, ErrRequestNotFound) {
			return nil, err
		}
		return nil, APIErrRequestIsScheduledToProcess
	}

	fileName := fileHeaders.Filename
	fileExt := filepath.Ext(fileName)

	if ok := s.validateUploadExtension(fileExt, ValidExstenstions); !ok {
		return nil, APIErrInvalidFileExstension
	}

	fileKey := s.generateUniqueKey(fileExt, fileExt)
	if err := s.uploader.Upload(ctx, fileKey, file); err != nil {
		return nil, err
	}

	req := NewRequest(userID, idempotencyKey, fileKey)

	if err := s.requestStore.Save(ctx, req); err != nil {
		return nil, err
	}

	return &CreateRequestResponse{
		ID:             req.ID(),
		IdempotencyKey: req.IdempotencyKey(),
		Status:         req.Status(),
		FileLink:       s.uploader.GenerateObjectURL(req.InputS3Key()),
		CreatedAt:      req.CreatedAt(),
	}, nil
}

func (s *service) generateUniqueKey(fileName string, ext string) string {
	return fmt.Sprintf("%v-%v%v", uuid.New(), fileName, ext)
}

func (s *service) validateUploadExtension(ext string, validExstensions []string) bool {
	return slices.Contains(validExstensions, ext)
}
