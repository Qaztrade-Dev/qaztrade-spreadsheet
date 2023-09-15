package service

import (
	"context"

	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
)

type SwitchStatusRequest struct {
	ApplicationID string
	StatusName    string
}

func (s *service) Revision(ctx context.Context, application *domain.Application) error {
	claims, err := authDomain.ExtractClaims[authDomain.UserClaims](ctx)
	if err != nil {
		return err
	}

	manager, err := s.mngRepo.GetCurrent(ctx, claims.UserID)
	if err != nil {
		return err
	}

	data, err := s.spreadsheetSvc.Comments(ctx, application)
	data.ManagerName = manager.Fullname
	data.ManagerEmail = manager.Email
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SwitchStatus(ctx context.Context, req *SwitchStatusRequest) error {
	application, err := s.applicationRepo.GetOne(ctx, &domain.GetManyInput{
		ApplicationID: req.ApplicationID,
	})
	if err != nil {
		return err
	}

	if err := s.applicationRepo.EditStatus(ctx, req.ApplicationID, req.StatusName); err != nil {
		return err
	}

	var (
		isManagerReviewing = (req.StatusName == domain.StatusManagerReviewing)
		isUserFilling      = (req.StatusName == domain.StatusUserFilling)
		isUserFixing       = (req.StatusName == domain.StatusUserFixing)

		mustSwitchModeEdit = false
		mustSwitchModeRead = false
	)

	switch {
	case isUserFilling:
		mustSwitchModeEdit = true
	case isUserFixing:
		mustSwitchModeEdit = true
		err := s.Revision(ctx, application)
		if err != nil {
			return err
		}
	case isManagerReviewing:
		mustSwitchModeRead = true
	default:
		mustSwitchModeRead = true
	}

	if mustSwitchModeEdit {
		if err := s.spreadsheetSvc.SwitchModeEdit(ctx, application.SpreadsheetID); err != nil {
			return err
		}
	}

	if mustSwitchModeRead {
		if err := s.spreadsheetSvc.SwitchModeRead(ctx, application.SpreadsheetID); err != nil {
			return err
		}
	}

	return nil
}
