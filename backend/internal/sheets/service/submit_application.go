package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
)

type SubmitApplicationRequest struct {
	SpreadsheetID string
	Application   *domain.Application
}

func (s *service) SubmitApplication(ctx context.Context, req *SubmitApplicationRequest) error {
	if err := s.sheetsRepo.UpdateApplication(ctx, req.SpreadsheetID, req.Application); err != nil {
		return err
	}

	return nil
}
