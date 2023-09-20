package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
)

type ForgotRequest struct {
	Email string
}

func (s *service) Forgot(ctx context.Context, input *ForgotRequest) error {
	email := domain.CleanEmail(input.Email)

	user, err := s.authRepo.GetOne(ctx, &domain.GetQuery{Email: email})
	if err != nil {
		return err
	}

	roles, err := s.authRepo.GetRoles(ctx, user.ID)
	if err != nil {
		return err
	}

	creds, err := s.credsRepo.Create(ctx, &domain.UserClaims{
		UserID: user.ID,
		Roles:  roles,
	})
	if err != nil {
		return err
	}

	if err := s.emailSvc.Send(ctx, input.Email, mailName, &MailPayload{
		Credentials: creds,
	}); err != nil {
		return err
	}

	return nil
}

const mailName = "forgot"

type MailPayload struct {
	Credentials *domain.Credentials
}
