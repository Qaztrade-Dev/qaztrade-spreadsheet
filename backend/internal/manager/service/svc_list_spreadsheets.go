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
	SignedAtFrom     time.Time
	SignedAtUntil    time.Time
}

func (s *service) ListSpreadsheets(ctx context.Context, req *ListSpreadsheetsRequest) (*domain.ApplicationList, error) {
	list, err := s.applicationRepo.GetMany(ctx, &domain.ApplicationQuery{
		Limit:            req.Limit,
		Offset:           req.Offset,
		BIN:              req.BIN,
		CompensationType: req.CompensationType,
		SignedAtFrom:     req.SignedAtFrom,
		SignedAtUntil:    req.SignedAtUntil,
	})
	if err != nil {
		return nil, err
	}

	return list, err
}
