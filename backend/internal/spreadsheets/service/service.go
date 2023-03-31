package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
)

type Service interface {
	CreateSpreadsheet(ctx context.Context, req *CreateSpreadsheetRequest) (spreadsheetID string, err error)
}

type service struct {
	spreadsheetSvc  domain.SpreadsheetService
	applicationRepo domain.ApplicationRepository
	userRepo        domain.UserRepository
}

func NewService() Service {
	return &service{}
}
