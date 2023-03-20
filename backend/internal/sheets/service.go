package sheets

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"github.com/pkg/errors"
)

type Service interface {
	SubmitRecord(ctx context.Context, req *SubmitRecordRequest) error
}

type SubmitRecordRequest struct {
	SpreadsheetID string
	Payload       *domain.Payload
}

type service struct {
	sheetsRepo domain.SheetsRepository
}

func NewService(sheetsRepo domain.SheetsRepository) Service {
	return &service{
		sheetsRepo: sheetsRepo,
	}
}

func (s *service) SubmitRecord(ctx context.Context, req *SubmitRecordRequest) error {
	if err := req.validate(); err != nil {
		return err
	}

	if err := s.sheetsRepo.InsertRecord(ctx, req.SpreadsheetID, req.Payload); err != nil {
		return err
	}

	return nil
}

var (
	ErrorSpreadsheetID = errors.New("empty SpreadsheetID")
	ErrorPayload       = errors.New("empty Payload")
)

func (r *SubmitRecordRequest) validate() error {
	if r.SpreadsheetID == "" {
		return ErrorSpreadsheetID
	}

	if r.Payload == nil {
		return ErrorPayload
	}

	return nil
}
