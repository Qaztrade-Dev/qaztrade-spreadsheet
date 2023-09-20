package adapters

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	spreadsheetsAdapters "github.com/doodocs/qaztrade/backend/internal/spreadsheets/adapters"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func (s *SpreadsheetServiceGoogle) AddTotalSumCells(ctx context.Context) error {
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
				totalSum_Затраты_на_доставку_транспортом(batch, sheetID)
			case "Затраты на сертификацию предприятия",
				"Затраты на рекламу ИКУ за рубежом",
				"Затраты на перевод каталога ИКУ",
				"Затраты на аренду помещения ИКУ",
				"Затраты на сертификацию ИКУ",
				"Затраты на демонстрацию ИКУ",
				"Затраты на франчайзинг",
				"Затраты на регистрацию товарных знаков",
				"Затраты на аренду",
				"Затраты на перевод",
				"Затраты на рекламу товаров за рубежом",
				"Затраты на участие в выставках",
				"Затраты на участие в выставках ИКУ",
				"Затраты на соответствие товаров требованиям":
				totalSum_Затраты_на_сертификацию_предприятия(batch, sheetID)
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

func totalSum(sheetID int64, parentA1, targetA1 string) []*sheets.Request {
	var (
		result = make([]*sheets.Request, 0)

		a1Range = A1ToRange(parentA1)
		cell    = A1ToCell(targetA1)
	)

	beforeTargetRange := &Range{
		From: a1Range.From,
		To: &Cell{
			Col: cell.Col - 1,
			Row: cell.Row,
		},
	}

	afterTargetRange := &Range{
		From: &Cell{
			Col: cell.Col + 1,
			Row: cell.Row,
		},
		To: a1Range.To,
	}

	result = append(result,
		UnmergeRequest(sheetID, parentA1),
		SetCellFormula(sheetID, targetA1, fmt.Sprintf("=sum(%[1]s4:%[1]s)", numberToColumn(cell.Col))),
	)

	// fmt.Println("--------")

	// fmt.Println(parentA1, targetA1)
	// fmt.Printf("=sum(%[1]s4:%[1]s)\n", numberToColumn(cell.Col))

	if beforeTargetRange.From.Col < cell.Col && !beforeTargetRange.From.Equals(beforeTargetRange.To) {
		result = append(result, MergeRequest(sheetID, beforeTargetRange.ToA1()))
		// fmt.Printf("merge %s\n", beforeTargetRange.ToA1())
	}

	if afterTargetRange.To.Col > cell.Col && !afterTargetRange.From.Equals(afterTargetRange.To) {
		result = append(result, MergeRequest(sheetID, afterTargetRange.ToA1()))
		// fmt.Printf("merge %s\n", afterTargetRange.ToA1())
	}

	// fmt.Println("--------")

	return result
}

// Затраты на доставку транспортом
func totalSum_Затраты_на_доставку_транспортом(batch *spreadsheetsAdapters.BatchUpdate, sheetID int64) {
	args := []struct {
		parentA1 string
		targetA1 string
	}{
		{"U2:Y2", "X2"},
		{"Z2:AF2", "AF2"},
		{"AG2:AL2", "AJ2"},
		{"AM2:AP2", "AP2"},
		{"AQ2:AT2", "AT2"},
		{"AU2:AX2", "AX2"},
	}

	for _, arg := range args {
		batch.WithRequest(
			totalSum(sheetID, arg.parentA1, arg.targetA1)...,
		)
	}
}

// Затраты на сертификацию предприятия
func totalSum_Затраты_на_сертификацию_предприятия(batch *spreadsheetsAdapters.BatchUpdate, sheetID int64) {
	args := []struct {
		parentA1 string
		targetA1 string
	}{
		{"N2:V2", "U2"},
		{"W2:AB2", "AA2"},
		{"AC2:AI2", "AG2"},
		{"AJ2:AN2", "AM2"},
		{"AO2:AR2", "AR2"},
		{"AS2:AV2", "AV2"},
	}

	for _, arg := range args {
		batch.WithRequest(
			totalSum(sheetID, arg.parentA1, arg.targetA1)...,
		)
	}
}

func SetCellFormula(sheetID int64, a1, formula string) *sheets.Request {
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

type Cell struct {
	Col int64
	Row int64
}

type Range struct {
	From *Cell
	To   *Cell
}

func (cell *Cell) Equals(b *Cell) bool {
	return cell.Col == b.Col && cell.Row == b.Row
}

func (cell *Cell) ToA1() string {
	return numberToColumn(cell.Col) + strconv.Itoa(int(cell.Row)+1)
}

func (r *Range) ToA1() string {
	return r.From.ToA1() + ":" + r.To.ToA1()
}

func A1ToRange(fromToA1 string) *Range {
	splitted := strings.Split(fromToA1, ":")
	if len(splitted) != 2 {
		return nil
	}

	var (
		fromCell = A1ToCell(splitted[0])
		toCell   = A1ToCell(splitted[1])
	)

	return &Range{
		From: fromCell,
		To:   toCell,
	}
}
