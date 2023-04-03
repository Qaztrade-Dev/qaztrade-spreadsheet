package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
)

type SwitchStatusRequest struct {
	ApplicationID string
	StatusName    string
}

func (s *service) SwitchStatus(ctx context.Context, req *SwitchStatusRequest) error {
	application, err := s.applicationRepo.GetOne(ctx, &domain.ApplicationQuery{
		ApplicationID: req.ApplicationID,
	})
	if err != nil {
		return err
	}

	if err := s.applicationRepo.EditStatus(ctx, req.ApplicationID, req.StatusName); err != nil {
		return err
	}

	var (
		isUserFilling = (req.StatusName == domain.StatusUserFilling)
		isUserFixing  = (req.StatusName == domain.StatusUserFixing)
	)

	if isUserFilling || isUserFixing {
		if err := s.spreadsheetSvc.SwitchModeEdit(ctx, application.SpreadsheetID); err != nil {
			return err
		}
	} else {
		if err := s.spreadsheetSvc.SwitchModeRead(ctx, application.SpreadsheetID); err != nil {
			return err
		}
	}

	return nil
}
