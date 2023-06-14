package service

import (
	"context"
	"fmt"

	"github.com/doodocs/qaztrade/backend/internal/sign/domain"
)

type SyncSpreadsheetsRequest struct {
	SpreadsheetID string
}

func (s *service) SyncSpreadsheets(ctx context.Context, req *SyncSpreadsheetsRequest) error {
	sheets, err := s.spreadsheetRepo.GetSheets(ctx, req.SpreadsheetID)
	if err != nil {
		return err
	}
	if len(sheets) == 0 {
		return domain.ErrorEmptySpreadsheet
	}

	application, err := s.spreadsheetRepo.GetApplication(ctx, req.SpreadsheetID)
	if err != nil {
		return err
	}

	var (
		aggSheetTitles   string  = domain.SheetTitlesJoined(sheets)
		aggSheetExpenses float64 = domain.SheetsTotalExpenses(sheets)
	)

	application.ExpensesList = aggSheetTitles
	application.ExpensesSum = fmt.Sprintf("%f", aggSheetExpenses)
	application.ApplicationDate = domain.GetApplicationDate()

	if err := s.applicationRepo.AssignAttrs(ctx, req.SpreadsheetID, &domain.ApplicationAttrs{
		Application: application,
		SheetsAgg:   domain.SheetsAgg(sheets),
		Sheets:      sheets,
	}); err != nil {
		return err
	}

	return nil
}
