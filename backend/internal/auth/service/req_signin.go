package service

import (
	"context"
	"fmt"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
)

type SignInRequest struct {
	Email    string
	Password string
}

func (s *service) SignIn(ctx context.Context, input *SignInRequest) (*domain.Credentials, error) {
	email := domain.CleanEmail(input.Email)

	user, err := s.authRepo.SignIn(ctx, &domain.SignInInput{
		Email:    email,
		Password: input.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("error SignIn: %w", err)
	}

	roles, err := s.authRepo.GetRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("error GetRoles: %w", err)
	}

	creds, err := s.credsRepo.Create(ctx, &domain.UserClaims{
		UserID: user.ID,
		Roles:  roles,
	})

	if err != nil {
		return nil, fmt.Errorf("error Create: %w", err)
	}

	return creds, nil
}
