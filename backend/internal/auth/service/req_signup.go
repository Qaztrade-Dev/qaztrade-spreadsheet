package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
)

type SignUpRequest struct {
	Email    string
	Password string
	OrgName  string
}

func (s *service) SignUp(ctx context.Context, req *SignUpRequest) (*domain.Credentials, error) {
	userID, err := s.authRepo.SignUp(ctx, &domain.SignUpInput{
		Email:    req.Email,
		Password: req.Password,
		OrgName:  req.OrgName,
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
