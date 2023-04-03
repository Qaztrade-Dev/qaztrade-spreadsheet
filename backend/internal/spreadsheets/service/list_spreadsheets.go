package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
)

type ListSpreadsheetsRequest struct {
	Limit  uint64
	Offset uint64
	UserID string
}

func (s *service) ListSpreadsheets(ctx context.Context, req *ListSpreadsheetsRequest) (*domain.ApplicationList, error) {
	list, err := s.applicationRepo.GetMany(ctx, &domain.ApplicationQuery{
		UserID: req.UserID,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		return nil, err
	}

	return list, err
}
