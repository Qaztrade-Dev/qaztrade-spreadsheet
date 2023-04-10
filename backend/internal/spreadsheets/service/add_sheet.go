package service

import (
	"context"
)

type AddSheetRequest struct {
	SpreadsheetID string
	SheetName     string
}

func (s *service) AddSheet(ctx context.Context, req *AddSheetRequest) error {
	if err := s.spreadsheetSvc.AddSheet(ctx, req.SpreadsheetID, req.SheetName); err != nil {
		return err
	}

	return nil
}
