package service

import (
	"context"
)

type SyncSigningTimeRequest struct {
	SignDocumentID string
}

func (s *service) SyncSigningTime(ctx context.Context, req *SyncSigningTimeRequest) error {
	signedAt, err := s.signSvc.GetSigningTime(ctx, req.SignDocumentID)
	if err != nil {
		return err
	}

	application, err := s.applicationRepo.GetApplicationByDocumentID(ctx, req.SignDocumentID)
	if err != nil {
		return err
	}

	if err := s.applicationRepo.ConfirmSigningInfo(ctx, application.SpreadsheetID, signedAt); err != nil {
		return err
	}

	return nil
}
