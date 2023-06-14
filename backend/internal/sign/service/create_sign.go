package service

import (
	"context"
	"fmt"

	"github.com/doodocs/qaztrade/backend/internal/sign/domain"
)

type CreateSignRequest struct {
	SpreadsheetID string
}

func (s *service) CreateSign(ctx context.Context, req *CreateSignRequest) (string, error) {
	var linkbase = "https://link.doodocs.kz/"

	signApplication, err := s.applicationRepo.GetApplication(ctx, req.SpreadsheetID)
	if err != nil {
		return "", err
	}

	if signApplication.Status != domain.StatusUserFilling {
		return linkbase + signApplication.SignLink, nil
	}

	sheets, err := s.spreadsheetRepo.GetSheets(ctx, req.SpreadsheetID)
	if err != nil {
		return "", err
	}
	if len(sheets) == 0 {
		return "", domain.ErrorEmptySpreadsheet
	}

	hasMergedCells, err := s.spreadsheetRepo.HasMergedCells(ctx, req.SpreadsheetID, sheets)
	if err != nil {
		return "", domain.ErrorAbsentExpenses
	}
	if hasMergedCells {
		return "", domain.ErrorSpreadsheetHasMergedCells
	}

	application, err := s.spreadsheetRepo.GetApplication(ctx, req.SpreadsheetID)
	if err != nil {
		return "", err
	}

	var (
		aggSheetTitles   string  = domain.SheetTitlesJoined(sheets)
		aggSheetExpenses float64 = domain.SheetsTotalExpenses(sheets)
		documentName             = domain.GetDocumentName(application.Bin)
	)

	if aggSheetExpenses == 0 {
		return "", domain.ErrorExpensesZero
	}

	application.ExpensesList = aggSheetTitles
	application.ExpensesSum = fmt.Sprintf("%f", aggSheetExpenses)
	application.ApplicationDate = domain.GetApplicationDate()

	attachments, err := s.spreadsheetRepo.GetAttachments(ctx, req.SpreadsheetID, sheets)
	if err != nil {
		return "", err
	}

	pdfToSign, err := s.pdfSvc.Create(application, attachments)
	if err != nil {
		return "", err
	}

	resp, err := s.signSvc.CreateSigningDocument(ctx, documentName, pdfToSign)
	if err != nil {
		return "", err
	}

	if err := s.applicationRepo.AssignSigningInfo(ctx, req.SpreadsheetID, resp); err != nil {
		return "", err
	}

	if err := s.applicationRepo.AssignAttrs(ctx, req.SpreadsheetID, &domain.ApplicationAttrs{
		Application: application,
		SheetsAgg:   domain.SheetsAgg(sheets),
		Sheets:      sheets,
	}); err != nil {
		return "", err
	}

	return linkbase + resp.SignLink, nil
}
