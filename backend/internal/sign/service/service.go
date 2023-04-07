package service

import (
	"context"
	"fmt"
	"os"

	"github.com/doodocs/qaztrade/backend/internal/sign/domain"
)

type Service interface {
	CreateSign(ctx context.Context, req *CreateSignRequest) error
}

type CreateSignRequest struct {
	SpreadsheetID string
}

func (s *service) CreateSign(ctx context.Context, req *CreateSignRequest) error {
	application, err := s.spreadsheetRepo.GetApplication(ctx, req.SpreadsheetID)
	if err != nil {
		return err
	}

	attachments, err := s.spreadsheetRepo.GetAttachments(ctx, req.SpreadsheetID)
	if err != nil {
		return err
	}

	pdfToSign, err := s.pdfSvc.Create(application, attachments)
	if err != nil {
		return err
	}

	fmt.Println(
		"WriteFile err", os.WriteFile("hello.pdf", pdfToSign.Bytes(), 0644),
	)

	return nil
}

type service struct {
	pdfSvc          domain.PDFService
	spreadsheetRepo domain.SpreadsheetRepository
}

func NewService(
	pdfSvc domain.PDFService,
	spreadsheetRepo domain.SpreadsheetRepository,
) Service {
	return &service{
		pdfSvc:          pdfSvc,
		spreadsheetRepo: spreadsheetRepo,
	}
}
