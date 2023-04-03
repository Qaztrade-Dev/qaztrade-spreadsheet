package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
)

type Service interface {
	CreateSpreadsheet(ctx context.Context, req *CreateSpreadsheetRequest) (publicLink string, err error)
	ListSpreadsheets(ctx context.Context, req *ListSpreadsheetsRequest) (*domain.ApplicationList, error)
}

type service struct {
	spreadsheetSvc  domain.SpreadsheetService
	applicationRepo domain.ApplicationRepository
	userRepo        domain.UserRepository
}

func NewService(
	spreadsheetSvc domain.SpreadsheetService,
	applicationRepo domain.ApplicationRepository,
	userRepo domain.UserRepository,
) Service {
	return &service{
		spreadsheetSvc:  spreadsheetSvc,
		applicationRepo: applicationRepo,
		userRepo:        userRepo,
	}
}
