package adapters

import (
	"context"
	"fmt"
	"strings"
	"time"

	spreadsheetsAdapters "github.com/doodocs/qaztrade/backend/internal/spreadsheets/adapters"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func (s *SpreadsheetServiceGoogle) AddСоответствиеЖДН(ctx context.Context) error {
	httpClient, err := s.oauth2.GetClient(ctx)
	if err != nil {
		return err
	}

	spreadsheetsSvc, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	// spreadsheetIDs, err := s.getSpreadsheets(ctx, httpClient)
	// if err != nil {
	// 	return err
	// }
	spreadsheetIDs := []string{"1oMJFttuiPxoBdejx3Ul3D2nscE0x45oeNfo-upVe1gE"}

	for i, spreadsheetID := range spreadsheetIDs {
		spreadsheetID := spreadsheetID
		fmt.Println(spreadsheetID)

		spreadsheet, err := spreadsheetsSvc.Spreadsheets.Get(spreadsheetID).Do()
		if err != nil {
			return err
		}

		batch := spreadsheetsAdapters.NewBatchUpdate(spreadsheetsSvc)

		for _, sheet := range spreadsheet.Sheets {
			sheetID := sheet.Properties.SheetId
			title := strings.ReplaceAll(sheet.Properties.Title, "⏳ (ожидайте) ", "")
			switch title {
			case "Затраты на доставку транспортом":
				СоответствиеЖДН(batch, sheetID)
			}
		}

		if err := batch.Do(ctx, spreadsheetID); err != nil {
			return err
		}

		fmt.Printf("%v/%v\n", i+1, len(spreadsheetIDs))
		time.Sleep(time.Second)
	}

	return nil
}

func СоответствиеЖДН(batch *spreadsheetsAdapters.BatchUpdate, sheetID int64) {
	batch.WithRequest(
		InsertColumnLeft(sheetID, "EF"),
		SetCellText(sheetID, "EF2", &SetCellTextInput{"Соответствие ЖДН (оцифровка)", true, 8}),
		MergeRequest(sheetID, "EF2:EF3"),
		SetDataValidationOneOf(sheetID, "EF4", []string{"да", "нет"}),
	)
}

func InsertColumnLeft(sheetID int64, columnA1 string) *sheets.Request {
	var (
		col = columnToNumber(columnA1)
	)

	return &sheets.Request{
		InsertDimension: &sheets.InsertDimensionRequest{
			Range: &sheets.DimensionRange{
				Dimension:  "COLUMNS",
				StartIndex: col - 1,
				EndIndex:   col,
				SheetId:    sheetID,
			},
			InheritFromBefore: true,
		},
	}
}

func SetDataValidationOneOf(sheetID int64, fromA1 string, oneOf []string) *sheets.Request {
	var (
		cell        = A1ToCell(fromA1)
		oneOfValues = make([]*sheets.ConditionValue, len(oneOf))
	)

	for i := range oneOf {
		oneOfValues[i] = &sheets.ConditionValue{UserEnteredValue: oneOf[i]}
	}

	return &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Range: &sheets.GridRange{
				StartRowIndex: cell.Row,
				// EndRowIndex:      1000000,
				StartColumnIndex: cell.Col,
				EndColumnIndex:   cell.Col + 1,
				SheetId:          sheetID,
			},
			Cell: &sheets.CellData{
				DataValidation: &sheets.DataValidationRule{
					Condition: &sheets.BooleanCondition{
						Type:   "ONE_OF_LIST",
						Values: oneOfValues,
					},
					ShowCustomUi: true,
				},
			},
			Fields: "dataValidation",
		},
	}
}

type SetCellTextInput struct {
	Text     string
	Bold     bool
	FontSize int64
}

func SetCellText(sheetID int64, a1 string, input *SetCellTextInput) *sheets.Request {
	cell := A1ToCell(a1)

	return &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetID,
				StartRowIndex:    cell.Row,
				EndRowIndex:      cell.Row + 1,
				StartColumnIndex: cell.Col,
				EndColumnIndex:   cell.Col + 1,
			},
			Cell: &sheets.CellData{
				UserEnteredValue: &sheets.ExtendedValue{
					StringValue: &input.Text,
				},
				UserEnteredFormat: &sheets.CellFormat{
					TextFormat: &sheets.TextFormat{
						Bold:     input.Bold,
						FontSize: input.FontSize,
					},
				},
			},
			Fields: "userEnteredValue,userEnteredFormat.textFormat.bold,userEnteredFormat.textFormat.fontSize",
		},
	}
}
