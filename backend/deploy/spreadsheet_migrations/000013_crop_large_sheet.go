package adapters

import (
	"context"
	"fmt"

	spreadsheetsAdapters "github.com/doodocs/qaztrade/backend/internal/spreadsheets/adapters"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func (s *SpreadsheetServiceGoogle) CropLargeSheet(ctx context.Context) error {
	httpClient, err := s.oauth2.GetClient(ctx)
	if err != nil {
		return err
	}

	spreadsheetsSvc, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	var (
		batch = spreadsheetsAdapters.NewBatchUpdate(spreadsheetsSvc)
	)

	batch.WithRequest(
		&sheets.Request{
			DeleteDimension: &sheets.DeleteDimensionRequest{
				Range: &sheets.DimensionRange{
					SheetId:    1366326786,
					Dimension:  "ROWS",
					StartIndex: 200,
				},
			},
		},
		emptyCellRequst("DX67"),
		emptyCellRequst("DZ67"),
		emptyCellRequst("EB67"),
		emptyCellRequst("EC67"),
		emptyCellRequst("CE67"),
		emptyCellRequst("Y67"),
	)

	spreadsheetID := "1AfhwOrdMQMhmbgJNQhNOgt50fAmN1hl65WVFSmXqS68"
	if err := batch.Do(ctx, spreadsheetID); err != nil {
		fmt.Println(err)
	}

	return nil
}

func emptyCellRequst(a1 string) *sheets.Request {
	cell := A1ToCell(a1)

	return &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Range: &sheets.GridRange{
				SheetId:          1366326786,
				StartRowIndex:    cell.Row,
				EndRowIndex:      cell.Row + 1,
				StartColumnIndex: cell.Col,
				EndColumnIndex:   cell.Col + 1,
			},
			Cell: &sheets.CellData{
				UserEnteredValue: &sheets.ExtendedValue{
					FormulaValue: nil,
				},
			},
			Fields: "userEnteredValue",
		},
	}
}
