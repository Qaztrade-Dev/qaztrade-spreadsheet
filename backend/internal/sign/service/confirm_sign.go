package service

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sign/domain"
)

type ConfirmSignRequest struct {
	SignDocumentID string
}

func (s *service) ConfirmSign(ctx context.Context, req *ConfirmSignRequest) error {
	if err := s.publisher.Publish(ctx, []byte(req.SignDocumentID)); err != nil {
		return err
	}

	signedAt, err := s.signSvc.GetSigningTime(ctx, req.SignDocumentID)
	if err != nil {
		return err
	}

	application, err := s.applicationRepo.GetApplicationByDocumentID(ctx, req.SignDocumentID)
	if err != nil {
		return err
	}

	if err := s.spreadsheetRepo.UpdateSigningTime(ctx, application.SpreadsheetID, signedAt); err != nil {
		return err
	}

	if err := s.applicationRepo.ConfirmSigningInfo(ctx, application.SpreadsheetID, signedAt); err != nil {
		return err
	}

	if err := s.applicationRepo.EditStatus(ctx, application.SpreadsheetID, domain.StatusManagerReviewing); err != nil {
		return err
	}

	if err := s.spreadsheetRepo.SwitchModeRead(ctx, application.SpreadsheetID); err != nil {
		return err
	}

	return nil
}
