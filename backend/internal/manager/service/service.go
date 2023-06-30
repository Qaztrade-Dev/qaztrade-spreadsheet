package service

import (
	"context"
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
)

type Service interface {
	SwitchStatus(ctx context.Context, req *SwitchStatusRequest) error
	ListSpreadsheets(ctx context.Context, req *ListSpreadsheetsRequest) (*domain.ApplicationList, error)
	DownloadArchive(ctx context.Context, req *DownloadArchiveRequest) (*DownloadArchiveResponse, error)
	GetDDCard(ctx context.Context, req *GetDDCardRequest) (*http.Response, error)

	// GetManagers returns a list of managers
	GetManagers(ctx context.Context) ([]*domain.Manager, error)
}

type service struct {
	spreadsheetSvc     domain.SpreadsheetService
	applicationRepo    domain.ApplicationRepository
	spreadsheetStorage domain.Storage
	signingSvc         domain.SigningService
	mngRepo            domain.ManagersRepository
}

func NewService(
	spreadsheetSvc domain.SpreadsheetService,
	applicationRepo domain.ApplicationRepository,
	spreadsheetStorage domain.Storage,
	signingSvc domain.SigningService,
	mngRepo domain.ManagersRepository,
) Service {
	return &service{
		spreadsheetSvc:     spreadsheetSvc,
		applicationRepo:    applicationRepo,
		spreadsheetStorage: spreadsheetStorage,
		signingSvc:         signingSvc,
		mngRepo:            mngRepo,
	}
}
