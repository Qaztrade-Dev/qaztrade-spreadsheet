package service

import (
	"context"
	"errors"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
)

type AddRowsRequest struct {
	SpreadsheetID string
	Input         *domain.AddRowsInput
}

func (s *service) AddRows(ctx context.Context, req *AddRowsRequest) error {
	switch req.Input.SheetName {
	case "Заявление", "ТНВЭД", "ОКВЭД":
		return errors.New("list not allowed")
	}

	if err := s.sheetsRepo.AddRows(ctx, req.SpreadsheetID, req.Input); err != nil {
		return err
	}

	return nil
}
