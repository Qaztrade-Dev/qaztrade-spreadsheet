package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
)

type Service interface {
	CreateSpreadsheet(ctx context.Context, req *CreateSpreadsheetRequest) (publicLink string, err error)
	ListSpreadsheets(ctx context.Context, req *ListSpreadsheetsRequest) (*domain.ApplicationList, error)
	AddSheet(ctx context.Context, req *AddSheetRequest) error
}

type service struct {
	spreadsheetSvc  domain.SpreadsheetService
	applicationRepo domain.ApplicationRepository
	userRepo        domain.UserRepository
	jwtcli          *jwt.Client
}

func NewService(
	spreadsheetSvc domain.SpreadsheetService,
	applicationRepo domain.ApplicationRepository,
	userRepo domain.UserRepository,
	jwtcli *jwt.Client,
) Service {
	return &service{
		spreadsheetSvc:  spreadsheetSvc,
		applicationRepo: applicationRepo,
		userRepo:        userRepo,
		jwtcli:          jwtcli,
	}
}
