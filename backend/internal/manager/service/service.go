package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
)

type Service interface {
	SwitchStatus(ctx context.Context, req *SwitchStatusRequest) error
	ListSpreadsheets(ctx context.Context, req *ListSpreadsheetsRequest) (*domain.ApplicationList, error)
	DownloadArchive(ctx context.Context, req *DownloadArchiveRequest) (*DownloadArchiveResponse, error)
}

type service struct {
	spreadsheetSvc     domain.SpreadsheetService
	applicationRepo    domain.ApplicationRepository
	spreadsheetStorage domain.Storage
}

func NewService(
	spreadsheetSvc domain.SpreadsheetService,
	applicationRepo domain.ApplicationRepository,
	spreadsheetStorage domain.Storage,
) Service {
	return &service{
		spreadsheetSvc:     spreadsheetSvc,
		applicationRepo:    applicationRepo,
		spreadsheetStorage: spreadsheetStorage,
	}
}
