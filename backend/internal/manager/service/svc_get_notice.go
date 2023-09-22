package service

import (
	"bytes"
	"context"

	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
)

type GetNoticeRequest struct {
	ApplicationID string
}

func (s *service) Revision(ctx context.Context, application *domain.Application) (*domain.Revision, error) {
	claims, err := authDomain.ExtractClaims[authDomain.UserClaims](ctx)
	if err != nil {
		return nil, err
	}

	manager, err := s.mngRepo.GetCurrent(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	data, err := s.spreadsheetSvc.Comments(ctx, application)
	if err != nil {
		return nil, err
	}

	data.ManagerName = manager.Fullname
	data.ManagerEmail = manager.Email

	return data, nil
}

func (s *service) GetNotice(ctx context.Context, req *GetNoticeRequest) (*bytes.Buffer, error) {
	application, err := s.applicationRepo.GetOne(ctx, &domain.GetManyInput{
		ApplicationID: req.ApplicationID,
	})
	if err != nil {
		return nil, err
	}

	if application.Status != domain.StatusManagerReviewing {
		return nil, domain.ErrorApplicationNotUnderReview
	}

	remarks, err := s.Revision(ctx, application)
	if err != nil {
		return nil, err
	}

	result, err := s.noticeSvc.Create(remarks)
	if err != nil {
		return nil, err
	}

	return result, nil
}
