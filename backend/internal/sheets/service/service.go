package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
)

type Service interface {
	SubmitRecord(ctx context.Context, req *SubmitRecordRequest) error
	SubmitApplication(ctx context.Context, req *SubmitApplicationRequest) error
	AddSheet(ctx context.Context, req *AddSheetRequest) error
}

type service struct {
	sheetsRepo domain.SheetsRepository
	storage    domain.Storage
}

func NewService(sheetsRepo domain.SheetsRepository) Service {
	return &service{
		sheetsRepo: sheetsRepo,
	}
}
