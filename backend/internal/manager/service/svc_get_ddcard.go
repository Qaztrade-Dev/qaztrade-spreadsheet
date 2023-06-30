package service

import (
	"context"
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
)

type GetDDCardRequest struct {
	ApplicationID string
}

func (s *service) GetDDCard(ctx context.Context, req *GetDDCardRequest) (*http.Response, error) {
	application, err := s.applicationRepo.GetOne(ctx, &domain.ApplicationQuery{
		ApplicationID: req.ApplicationID,
	})
	if err != nil {
		return nil, err
	}

	if application.Status == domain.StatusUserFilling {
		return nil, domain.ErrorApplicationNotSigned
	}

	httpResp, err := s.signingSvc.GetDDCard(ctx, application.SignDocumentID)
	if err != nil {
		return nil, err
	}

	return httpResp, nil
}
