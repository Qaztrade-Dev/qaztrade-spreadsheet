package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
)

type SubmitRecordRequest struct {
	SpreadsheetID string
	Payload       *domain.Payload
}

func (s *service) SubmitRecord(ctx context.Context, req *SubmitRecordRequest) error {
	if err := s.sheetsRepo.InsertRecord(ctx, req.SpreadsheetID, req.Payload); err != nil {
		return err
	}

	return nil
}
