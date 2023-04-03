package service

import (
	"context"
)

func (s *service) UpdateToken(ctx context.Context, authCode string) error {
	tok, err := s.config.Exchange(context.TODO(), authCode)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateToken(ctx, tok); err != nil {
		return err
	}

	return nil
}
