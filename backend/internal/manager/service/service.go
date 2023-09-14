package service

import (
	"bytes"
	"context"
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
)

type Service interface {
	SwitchStatus(ctx context.Context, req *SwitchStatusRequest) error
	ListSpreadsheets(ctx context.Context, req *ListSpreadsheetsRequest) (*domain.ApplicationList, error)
	GetDDCard(ctx context.Context, req *GetDDCardRequest) (*http.Response, error)

	// GetManagers returns a list of managers
	GetManagers(ctx context.Context) ([]*domain.Manager, error)
	GetNotice(ctx context.Context, req *GetNoticeRequest) (*bytes.Buffer, error)
	SendNotice(ctx context.Context, req *SendNoticeRequest) (string, error)
}

type service struct {
	spreadsheetSvc  domain.SpreadsheetService
	applicationRepo domain.ApplicationRepository
	signingSvc      domain.SigningService
	mngRepo         domain.ManagersRepository
	noticeSvc       domain.NoticeService
	storage         domain.Storage
	emailSvc        domain.EmailService
}

func NewService(
	spreadsheetSvc domain.SpreadsheetService,
	applicationRepo domain.ApplicationRepository,
	signingSvc domain.SigningService,
	mngRepo domain.ManagersRepository,
	noticeSvc domain.NoticeService,
	storage domain.Storage,
	emailSvc domain.EmailService,
) Service {
	return &service{
		spreadsheetSvc:  spreadsheetSvc,
		applicationRepo: applicationRepo,
		signingSvc:      signingSvc,
		mngRepo:         mngRepo,
		noticeSvc:       noticeSvc,
		storage:         storage,
		emailSvc:        emailSvc,
	}
}
