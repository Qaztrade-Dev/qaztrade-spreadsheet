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

func (s *service) ListSpreadsheets(ctx context.Context, input *domain.GetManyInput) (*domain.ApplicationList, error) {
	list, err := s.applicationRepo.GetMany(ctx, input)
	if err != nil {
		return nil, err
	}

	return list, err
}
