package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
)

type SignInRequest struct {
	Email    string
	Password string
}

func (s *service) SignIn(ctx context.Context, req *SignInRequest) (*domain.Credentials, error) {
	userID, err := s.authRepo.SignIn(ctx, &domain.SignInInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	creds, err := s.credsRepo.Create(ctx, userID)
	if err != nil {
		return nil, err
	}

	return creds, nil
}
