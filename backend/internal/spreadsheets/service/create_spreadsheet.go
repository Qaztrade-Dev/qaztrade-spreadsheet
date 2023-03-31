package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
)

type CreateSpreadsheetRequest struct {
	UserID string
}

func (s *service) CreateSpreadsheet(ctx context.Context, req *CreateSpreadsheetRequest) (string, error) {
	user, err := s.userRepo.Get(ctx, req.UserID)
	if err != nil {
		return "", err
	}

	spreadsheetID, err := s.spreadsheetSvc.Create(ctx, user)
	if err != nil {
		return "", err
	}

	publicLink := s.spreadsheetSvc.GetPublicLink(ctx, spreadsheetID)

	if err := s.applicationRepo.Create(ctx, req.UserID, &domain.Application{
		UserID:        req.UserID,
		SpreadsheetID: spreadsheetID,
		Link:          publicLink,
	}); err != nil {
		return "", err
	}

	return publicLink, nil
}
