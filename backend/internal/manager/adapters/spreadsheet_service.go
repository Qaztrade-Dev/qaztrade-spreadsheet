package adapters

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

type SpreadsheetServiceGoogle struct {
	service *drive.Service
}

var _ domain.SpreadsheetService = (*SpreadsheetServiceGoogle)(nil)

func NewSpreadsheetService(ctx context.Context, credentialsJson []byte) (*SpreadsheetServiceGoogle, error) {
	service, err := drive.NewService(ctx, option.WithCredentialsJSON(credentialsJson))
	if err != nil {
		return nil, err
	}

	return &SpreadsheetServiceGoogle{
		service: service,
	}, err
}

func (s *SpreadsheetServiceGoogle) SwitchModeRead(ctx context.Context, spreadsheetID string) error {
	permission := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}
	_, err := s.service.Permissions.Insert(spreadsheetID, permission).Do()
	if err != nil {
		return err
	}
	return nil
}

func (s *SpreadsheetServiceGoogle) SwitchModeEdit(ctx context.Context, spreadsheetID string) error {
	permission := &drive.Permission{
		Type: "anyone",
		Role: "writer",
	}
	_, err := s.service.Permissions.Insert(spreadsheetID, permission).Do()
	if err != nil {
		return err
	}
	return nil
}
