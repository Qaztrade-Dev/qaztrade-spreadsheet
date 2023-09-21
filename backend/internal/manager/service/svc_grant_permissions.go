package service

import (
	"context"

	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
)

type GrantPermissionsRequest struct {
	ApplicationID string
}

func (s *service) GrantPermissions(ctx context.Context, req *GrantPermissionsRequest) error {
	claims, err := authDomain.ExtractClaims[authDomain.UserClaims](ctx)
	if err != nil {
		return err
	}

	manager, err := s.mngRepo.GetCurrent(ctx, claims.UserID)
	if err != nil {
		return err
	}

	application, err := s.applicationRepo.GetOne(ctx, &domain.GetManyInput{
		ApplicationID: req.ApplicationID,
	})
	if err != nil {
		return err
	}

	if err := s.spreadsheetSvc.GrantAdminPermissions(ctx, application.SpreadsheetID, manager.Email); err != nil {
		return err
	}

	return nil
}
