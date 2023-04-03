package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
)

type ForgotRequest struct {
	Email string
}

func (s *service) Forgot(ctx context.Context, req *ForgotRequest) error {
	user, err := s.authRepo.GetOne(ctx, &domain.GetQuery{
		Email: req.Email,
	})
	if err != nil {
		return err
	}

	creds, err := s.credsRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	if err := s.emailSvc.Send(ctx, req.Email, mailName, &MailPayload{
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
