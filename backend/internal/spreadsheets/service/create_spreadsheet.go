package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
)

type CreateSpreadsheetRequest struct {
	UserID string
}

func (s *service) CreateSpreadsheet(ctx context.Context, req *CreateSpreadsheetRequest) (string, error) {
	user, err := s.userRepo.Get(ctx, req.UserID)
	if err != nil {
		return "", err
	}

	spreadsheetID, err := s.spreadsheetSvc.Create(ctx, user)
	if err != nil {
		return "", err
	}

	publicLink := s.spreadsheetSvc.GetPublicLink(ctx, spreadsheetID)

	if err := s.applicationRepo.Create(ctx, req.UserID, &domain.Application{
		UserID:        req.UserID,
		SpreadsheetID: spreadsheetID,
		Link:          publicLink,
	}); err != nil {
		return "", err
	}

	tallyLink, err := s.getTallyLink(spreadsheetID)
	if err != nil {
		return "", err
	}

	return tallyLink, nil
}

func (s *service) getTallyLink(spreadsheetID string) (string, error) {
	baseURL := "https://tally.so/r/m6LxZB"
	tokenStr, err := jwt.NewTokenString(s.jwtcli, &domain.SpreadsheetClaims{
		SpreadsheetID: spreadsheetID,
	})
	if err != nil {
		return "", err
	}

	queryParams := url.Values{
		"token":          {tokenStr},
		"spreadsheet_id": {spreadsheetID},
	}

	return fmt.Sprintf("%s?%s", baseURL, queryParams.Encode()), nil
}
