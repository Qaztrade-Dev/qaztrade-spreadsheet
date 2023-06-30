package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
)

func (s *service) GetManagers(ctx context.Context) ([]*domain.Manager, error) {
	managers, err := s.mngRepo.GetMany(ctx)
	if err != nil {
		return nil, err
	}

	return managers, nil
}
