package adapters

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/api/sheets/v4"
)

type BatchUpdate struct {
	service  *sheets.Service
	requests []*sheets.Request
}

func NewBatchUpdate(service *sheets.Service) *BatchUpdate {
	return &BatchUpdate{
		service:  service,
		requests: make([]*sheets.Request, 0),
	}
}

func (b *BatchUpdate) WithProtectedRange(sheetID int64, protectedRanges []*sheets.ProtectedRange) {
	for _, pr := range protectedRanges {
		b.requests = append(b.requests, &sheets.Request{
			AddProtectedRange: &sheets.AddProtectedRangeRequest{
				ProtectedRange: &sheets.ProtectedRange{
					Range: &sheets.GridRange{
						SheetId:          sheetID,
						StartRowIndex:    pr.Range.StartRowIndex,
						EndRowIndex:      pr.Range.EndRowIndex,
						StartColumnIndex: pr.Range.StartColumnIndex,
						EndColumnIndex:   pr.Range.EndColumnIndex,
					},
					ProtectedRangeId:      pr.ProtectedRangeId,
					Description:           pr.Description,
					WarningOnly:           pr.WarningOnly,
					RequestingUserCanEdit: pr.RequestingUserCanEdit,
					Editors:               pr.Editors,
				},
			},
		})
	}
}

func (b *BatchUpdate) WithSheetName(sheetID int64, sheetName string) {
	b.requests = append(b.requests, &sheets.Request{
		UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
			Properties: &sheets.SheetProperties{
				SheetId: sheetID,
				Title:   sheetName,
			},
			Fields: "title",
		},
	})
}

func (b *BatchUpdate) WithRequest(requests ...*sheets.Request) {
	b.requests = append(b.requests, requests...)
}

func (b *BatchUpdate) Do(ctx context.Context, spreadsheetID string) error {
	if len(b.requests) == 0 {
		return nil
	}
	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{Requests: b.requests}
	_, err := b.service.Spreadsheets.BatchUpdate(spreadsheetID, batchUpdateRequest).Context(ctx).Do()
	return err
}

func UnmergeRequest(sheetID int64, fromToA1 string) *sheets.Request {
	splitted := strings.Split(fromToA1, ":")
	if len(splitted) != 2 {
		return nil
	}

	var (
		fromCol, fromRow = A1ToNumbers(splitted[0])
		toCol, toRow     = A1ToNumbers(splitted[1])
	)

	return &sheets.Request{
		UnmergeCells: &sheets.UnmergeCellsRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetID,
				StartRowIndex:    fromRow,
				EndRowIndex:      toRow + 1,
				StartColumnIndex: fromCol,
				EndColumnIndex:   toCol + 1,
			},
		},
	}
}

func MergeRequest(sheetID int64, fromToA1 string) *sheets.Request {
	splitted := strings.Split(fromToA1, ":")
	if len(splitted) != 2 {
		return nil
	}

	var (
		fromCol, fromRow = A1ToNumbers(splitted[0])
		toCol, toRow     = A1ToNumbers(splitted[1])
	)

	return &sheets.Request{
		MergeCells: &sheets.MergeCellsRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetID,
				StartRowIndex:    fromRow,
				EndRowIndex:      toRow + 1,
				StartColumnIndex: fromCol,
				EndColumnIndex:   toCol + 1,
			},
			MergeType: "MERGE_ALL",
		},
	}
}

type SetCellTextInput struct {
	Text     string
	Bold     bool
	FontSize int64
}

func SetCellText(sheetID int64, a1 string, input *SetCellTextInput) *sheets.Request {
	var (
		col, row = A1ToNumbers(a1)
	)

	return &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetID,
				StartRowIndex:    row,
				EndRowIndex:      row + 1,
				StartColumnIndex: col,
				EndColumnIndex:   col + 1,
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

func SetCellFormula(sheetID int64, a1, formula string) *sheets.Request {
	var (
		col, row = A1ToNumbers(a1)
	)

	return &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetID,
				StartRowIndex:    row,
				EndRowIndex:      row + 1,
				StartColumnIndex: col,
				EndColumnIndex:   col + 1,
			},
			Cell: &sheets.CellData{
				UserEnteredValue: &sheets.ExtendedValue{
					FormulaValue: &formula,
				},
				UserEnteredFormat: &sheets.CellFormat{
					NumberFormat: &sheets.NumberFormat{
						Type:    "NUMBER",
						Pattern: "#,##0.00",
					},
				},
			},
			Fields: "userEnteredValue,userEnteredFormat.numberFormat",
		},
	}
}

func columnToNumber(column string) int64 {
	column = strings.ToUpper(column)
	var num int64 = 0
	for _, char := range column {
		num = num*26 + int64(char-'A') + 1
	}
	return num
}

func A1ToNumbers(a1 string) (int64, int64) {
	for i, r := range a1 {
		if r >= '0' && r <= '9' {
			col := a1[:i]
			row, _ := strconv.Atoi(a1[i:])
			fmt.Println(a1, columnToNumber(col)-1, int64(row-1))
			return columnToNumber(col) - 1, int64(row - 1)
		}
	}
	return -1, -1 // Invalid A1 notation
}
