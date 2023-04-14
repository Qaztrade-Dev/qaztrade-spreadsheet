package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sign/domain"
)

type Service interface {
	CreateSign(ctx context.Context, req *CreateSignRequest) (string, error)
}

type service struct {
	pdfSvc          domain.PDFService
	signSvc         domain.SigningService
	spreadsheetRepo domain.SpreadsheetRepository
	applicationRepo domain.ApplicationRepository
}

func NewService(
	pdfSvc domain.PDFService,
	signSvc domain.SigningService,
	spreadsheetRepo domain.SpreadsheetRepository,
	applicationRepo domain.ApplicationRepository,
) Service {
	return &service{
		pdfSvc:          pdfSvc,
		signSvc:         signSvc,
		spreadsheetRepo: spreadsheetRepo,
		applicationRepo: applicationRepo,
	}
}
