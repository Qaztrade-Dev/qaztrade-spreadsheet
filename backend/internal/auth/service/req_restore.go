package service

import (
	"context"
)

type RestoreRequest struct {
	UserID   string
	Password string
}

func (s *service) Restore(ctx context.Context, req *RestoreRequest) error {
	if err := s.authRepo.UpdatePassword(ctx, req.UserID, req.Password); err != nil {
		return err
	}

	return nil
}
