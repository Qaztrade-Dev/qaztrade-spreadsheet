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
		isManagerReviewing = (req.StatusName == domain.StatusManagerReviewing)
		isUserFilling      = (req.StatusName == domain.StatusUserFilling)
		isUserFixing       = (req.StatusName == domain.StatusUserFixing)

		mustSwitchModeEdit        = false
		mustSwitchModeRead        = false
		mustBlockImportantRanges  = false
		mustUnlockImportantRanges = false
	)

	switch {
	case isUserFilling:
		mustSwitchModeEdit = true
		mustUnlockImportantRanges = true
	case isUserFixing:
		mustSwitchModeEdit = true
		mustUnlockImportantRanges = true
		_, err := s.Revision(ctx, application)
		if err != nil {
			return err
		}
	case isManagerReviewing:
		mustSwitchModeRead = true
		mustBlockImportantRanges = true
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

	if mustBlockImportantRanges {
		if err := s.spreadsheetSvc.BlockImportantRanges(ctx, application.SpreadsheetID); err != nil {
			return err
		}
	}

	if mustUnlockImportantRanges {
		if err := s.spreadsheetSvc.UnlockImportantRanges(ctx, application.SpreadsheetID); err != nil {
			return err
		}
	}

	return nil
}
