package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
)

type Service interface {
	SubmitApplication(ctx context.Context, req *SubmitApplicationRequest) error
	UploadFile(ctx context.Context, req *UploadFileRequest) error
	AddRows(ctx context.Context, req *AddRowsRequest) error
}

type service struct {
	sheetsRepo domain.SheetsRepository
	storage    domain.Storage
}

func NewService(sheetsRepo domain.SheetsRepository, storage domain.Storage) Service {
	return &service{
		sheetsRepo: sheetsRepo,
		storage:    storage,
	}
}
