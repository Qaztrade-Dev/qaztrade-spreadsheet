package service

import (
	"context"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/sign/domain"
)

type ConfirmSignRequest struct {
	SignDocumentID string
}

func (s *service) ConfirmSign(ctx context.Context, req *ConfirmSignRequest) error {
	signingTime, err := s.signSvc.GetSigningTime(ctx, req.SignDocumentID)
	if err != nil {
		return err
	}

	application, err := s.applicationRepo.GetApplicationByDocumentID(ctx, req.SignDocumentID)
	if err != nil {
		return err
	}

	signingTimeStr, err := createSigningTimeStr(signingTime)
	if err != nil {
		return err
	}

	if err := s.spreadsheetRepo.UpdateSigningTime(ctx, application.SpreadsheetID, signingTimeStr); err != nil {
		return err
	}

	if err := s.applicationRepo.ConfirmSigningInfo(ctx, application.SpreadsheetID, signingTimeStr); err != nil {
		return err
	}

	if err := s.applicationRepo.EditStatus(ctx, application.SpreadsheetID, domain.StatusManagerReviewing); err != nil {
		return err
	}

	if err := s.spreadsheetRepo.SwitchModeRead(ctx, application.SpreadsheetID); err != nil {
		return err
	}

	if err := s.spreadsheetRepo.BlockImportantRanges(ctx, application.SpreadsheetID); err != nil {
		return err
	}

	return nil
}

func createSigningTimeStr(tm time.Time) (string, error) {
	location, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		return "", err
	}

	timeStr := tm.In(location).Format("02.01.2006")
	return timeStr, nil
}
