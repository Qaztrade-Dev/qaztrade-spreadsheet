package adapters

import (
	"context"
	"fmt"
	"time"

	spreadsheetsAdapters "github.com/doodocs/qaztrade/backend/internal/spreadsheets/adapters"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func (s *SpreadsheetServiceGoogle) FixUserFixingSpreadsheets(ctx context.Context) error {
	httpClient, err := s.oauth2.GetClient(ctx)
	if err != nil {
		return err
	}

	driveSvc, err := drive.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	spreadsheetsSvc, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	spreadsheetIDs := []string{
		"1dnxGao23OFWdT0v_ry5AcGn9Rrk7DQP62KO6Kg5Xymw",
		"1OdfMSikfzh6NvAq2rs1bFVWgs9PWVtUBCN4gCvZgW_I",
		// "15ez9Hs8SF0RMhcAPRnzs7ptcU0j0_fYesJLOvZA3MNo",
	}

	for i, spreadsheetID := range spreadsheetIDs {
		spreadsheetID := spreadsheetID
		fmt.Println(spreadsheetID)

		if err := SwitchModeEdit(ctx, driveSvc, spreadsheetID); err != nil {
			return err
		}

		if err := LockSheets(ctx, spreadsheetsSvc, spreadsheetID); err != nil {
			return err
		}

		fmt.Printf("%v/%v\n", i+1, len(spreadsheetIDs))
		time.Sleep(time.Second)
	}

	return nil
}

func SwitchModeEdit(ctx context.Context, driveSvc *drive.Service, spreadsheetID string) error {
	permission := &drive.Permission{
		Type: "anyone",
		Role: "writer",
	}

	_, err := driveSvc.Permissions.Create(spreadsheetID, permission).Do()
	if err != nil {
		return err
	}
	return nil
}

func LockSheets(ctx context.Context, sheetsSvc *sheets.Service, spreadsheetID string) error {
	spreadsheet, err := sheetsSvc.Spreadsheets.Get(spreadsheetID).Context(ctx).Do()
	if err != nil {
		return err
	}

	batch := spreadsheetsAdapters.NewBatchUpdate(sheetsSvc)

	for _, sheet := range spreadsheet.Sheets {
		sheet := sheet

		if sheet.Properties.Title == "Заявление" ||
			sheet.Properties.Title == "ТНВЭД" ||
			sheet.Properties.Title == "ОКВЭД" {
			continue
		}

		batch.WithRequest(
			&sheets.Request{
				AddProtectedRange: &sheets.AddProtectedRangeRequest{
					ProtectedRange: &sheets.ProtectedRange{
						Range: &sheets.GridRange{
							SheetId: sheet.Properties.SheetId,
						},
						Description: "Protecting entire sheet" + " " + sheet.Properties.Title,
						WarningOnly: false,
					},
				},
			},
		)
	}

	if err := batch.Do(ctx, spreadsheetID); err != nil {
		return err
	}

	return nil
}
