package service

import (
	"context"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
)

type ListSpreadsheetsRequest struct {
	Limit            uint64
	Offset           uint64
	BIN              string
	CompensationType string
	SignedAt         time.Time
}

func (s *service) ListSpreadsheets(ctx context.Context, req *ListSpreadsheetsRequest) (*domain.ApplicationList, error) {
	list, err := s.applicationRepo.GetMany(ctx, &domain.ApplicationQuery{
		Limit:            req.Limit,
		Offset:           req.Offset,
		BIN:              req.BIN,
		CompensationType: req.CompensationType,
		SignedAt:         req.SignedAt,
	})
	if err != nil {
		return nil, err
	}

	return list, err
}
