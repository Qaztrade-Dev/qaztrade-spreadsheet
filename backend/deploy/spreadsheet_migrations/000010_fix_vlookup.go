package adapters

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	spreadsheetsAdapters "github.com/doodocs/qaztrade/backend/internal/spreadsheets/adapters"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func (s *SpreadsheetServiceGoogle) FixVLOOKUP(ctx context.Context) error {
	httpClient, err := s.oauth2.GetClient(ctx)
	if err != nil {
		return err
	}

	spreadsheetsSvc, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return err
	}

	spreadsheetIDs, err := s.getSpreadsheets(ctx, httpClient)
	if err != nil {
		return err
	}

	// spreadsheetIDs := []string{"1UDWxrcOuIB-XyL9Dld1THLUpGZidBMHQBHSld-Tb8y8"}

	for i, spreadsheetID := range spreadsheetIDs {
		spreadsheetID := spreadsheetID
		fmt.Println(spreadsheetID)

		spreadsheet, err := spreadsheetsSvc.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Do()
		if err != nil {
			return err
		}

		var sheet *sheets.Sheet = nil
		for i := range spreadsheet.Sheets {
			title := strings.ReplaceAll(spreadsheet.Sheets[i].Properties.Title, "⏳ (ожидайте) ", "")
			if title == "Затраты на доставку транспортом" {
				sheet = spreadsheet.Sheets[i]
				break
			}
		}

		if sheet == nil {
			continue
		}

		// add check for DX is Код ТНВЭД
		strValue := sheet.Data[0].RowData[1].Values[127].UserEnteredValue.StringValue
		if strValue != nil && *strValue == "Код ТНВЭД" {
			fmt.Println("skip")
			continue
		}

		batch := spreadsheetsAdapters.NewBatchUpdate(spreadsheetsSvc)
		batch.WithRequest(
			InsertColumnLeft(sheet.Properties.SheetId, "DX"),
			SetCellText(sheet.Properties.SheetId, "DX2", &SetCellTextInput{"Код ТНВЭД", true, 8}),
			MergeRequest(sheet.Properties.SheetId, "DX2:DX3"),
			&sheets.Request{
				RepeatCell: &sheets.RepeatCellRequest{
					Range: &sheets.GridRange{
						SheetId:          sheet.Properties.SheetId,
						StartRowIndex:    3,
						StartColumnIndex: 127,
						EndColumnIndex:   128,
					},
					Fields: "*",
					Cell: &sheets.CellData{
						UserEnteredValue: &sheets.ExtendedValue{
							StringValue: nil,
						},
						UserEnteredFormat: &sheets.CellFormat{
							NumberFormat: &sheets.NumberFormat{
								Type: "TEXT",
							},
							BackgroundColor: &sheets.Color{
								Red:   243.0 / 255.0,
								Green: 243.0 / 255.0,
								Blue:  243.0 / 255.0,
								Alpha: 1.0,
							},
						},
						DataValidation: nil,
					},
				},
			},
			ClearColumn(sheet.Properties.SheetId, "DX"),
			ClearColumn(sheet.Properties.SheetId, "DZ"),
			ClearColumn(sheet.Properties.SheetId, "EB"),
			ClearColumn(sheet.Properties.SheetId, "EC"),
			SetCellArrayFormula(sheet.Properties.SheetId, "DX4", `=ARRAYFORMULA(IF(ISBLANK(DW4:DW); ""; REGEXEXTRACT(DW4:DW; "(.+?)\s-\s")))`),
			SetCellArrayFormula(sheet.Properties.SheetId, "DZ4", `=ARRAYFORMULA(IF(DY4:DY="нет"; "отсутствует"; IF(LEN(DX4:DX)=0; ""; VLOOKUP(DX4:DX; 'ТНВЭД'!C:E; 3; FALSE))))`),
			SetCellArrayFormula(sheet.Properties.SheetId, "EB4", `=ARRAYFORMULA(IF(((EA4:EA)="") + ((DZ4:DZ)=""); ""; SWITCH(DZ4:DZ; "высокий"; EA4:EA*0,8; "средний"; EA4:EA*0,5; "нижний"; EA4:EA*0,3; 0)))`),
			SetCellArrayFormula(sheet.Properties.SheetId, "EC4", `=ARRAYFORMULA(IF(((AL4:AL)="") + ((EB4:EB)=""); ""; SWITCH(AL4:AL; "да"; EB4:EB*1,05; "нет"; EB4:EB*1; 0)))`),
		)
		if err := batch.Do(ctx, spreadsheetID); err != nil {
			return err
		}

		fmt.Printf("%v/%v\n", i+1, len(spreadsheetIDs))
		time.Sleep(time.Second)
	}

	return nil
}

func SetCellArrayFormula(sheetID int64, a1, formula string) *sheets.Request {
	cell := A1ToCell(a1)

	return &sheets.Request{
		UpdateCells: &sheets.UpdateCellsRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetID,
				StartRowIndex:    cell.Row,
				EndRowIndex:      cell.Row + 1,
				StartColumnIndex: cell.Col,
				EndColumnIndex:   cell.Col + 1,
			},
			Fields: "userEnteredValue",
			Rows: []*sheets.RowData{
				{
					Values: []*sheets.CellData{
						{
							UserEnteredValue: &sheets.ExtendedValue{
								FormulaValue: aws.String(formula),
							},
						},
					},
				},
			},
		},
	}
}

func ClearColumn(sheetID int64, col string) *sheets.Request {
	cell := columnToNumber(col) - 1

	return &sheets.Request{
		UpdateCells: &sheets.UpdateCellsRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetID,
				StartRowIndex:    3,
				StartColumnIndex: cell,
				EndColumnIndex:   cell + 1,
			},
			Fields: "userEnteredValue",
			Rows: []*sheets.RowData{
				{
					Values: []*sheets.CellData{
						{
							UserEnteredValue: &sheets.ExtendedValue{
								StringValue:  nil,
								FormulaValue: nil,
								NumberValue:  nil,
								BoolValue:    nil,
							},
						},
					},
				},
			},
		},
	}
}
