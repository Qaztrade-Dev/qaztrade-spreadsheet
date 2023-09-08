package adapters

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	spreadsheetsAdapters "github.com/doodocs/qaztrade/backend/internal/spreadsheets/adapters"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func (s *SpreadsheetServiceGoogle) ConvertToFinancialSumCells(ctx context.Context) error {
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

	for i := 626; i < len(spreadsheetIDs); i++ {
		spreadsheetID := spreadsheetIDs[i]
		fmt.Println(spreadsheetID)

		spreadsheet, err := spreadsheetsSvc.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Do()
		if err != nil {
			return err
		}

		batch := spreadsheetsAdapters.NewBatchUpdate(spreadsheetsSvc)
		for _, sheet := range spreadsheet.Sheets {
			sheetID := sheet.Properties.SheetId
			title := strings.ReplaceAll(sheet.Properties.Title, "⏳ (ожидайте) ", "")

			if _, ok := sums[title]; ok {
				convertToFinancialSumCells_Expenses(batch, sheetID, sheet.Data[0], title)
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

var nonAlphanumericRegex = regexp.MustCompile(`[^0-9,]`)

func getColumnData(sheet *sheets.GridData, columnIdx, rowIdx int) [][]interface{} {
	data := make([][]interface{}, len(sheet.RowData[rowIdx:]))
	for i, row := range sheet.RowData[rowIdx:] {
		if columnIdx >= len(row.Values) {
			continue
		}

		cell := row.Values[columnIdx]

		if cell == nil || cell.UserEnteredValue == nil {
			continue
		}

		value := ""

		if cell.UserEnteredValue.StringValue != nil {
			value = *cell.UserEnteredValue.StringValue
			value = nonAlphanumericRegex.ReplaceAllString(value, "")
		} else if cell.UserEnteredValue.NumberValue != nil {
			value = fmt.Sprintf("%f", *cell.UserEnteredValue.NumberValue)
			value = strings.TrimRight(value, "0")
			if value[len(value)-1] == '.' {
				value = value[:len(value)-1]
			}
			value = strings.ReplaceAll(value, ".", ",")
		} else {
			value = ""
		}

		value = strings.ReplaceAll(value, ",", ".")
		floatValue, err := strconv.ParseFloat(value, 64)
		if err == nil {
			data[i] = append(data[i], floatValue)
		} else {
			data[i] = append(data[i], value)
		}
	}
	return data
}

func convertToFinancialSumCells(sheetID int64, targetColumn string, sheet *sheets.GridData, title string) ([]*sheets.Request, []*sheets.ValueRange) {
	var (
		requests    = make([]*sheets.Request, 0)
		valueRanges = make([]*sheets.ValueRange, 0)

		column        = columnToNumber(targetColumn) - 1
		fromRow int64 = 3
	)

	requests = append(requests,
		&sheets.Request{
			FindReplace: &sheets.FindReplaceRequest{
				Find:          "[^0-9,]",
				Replacement:   "",
				SearchByRegex: true,
				Range: &sheets.GridRange{
					SheetId:          sheetID,
					StartColumnIndex: column,
					EndColumnIndex:   column + 1,
					StartRowIndex:    fromRow,
				},
			},
		},
		&sheets.Request{
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
					SheetId:          sheetID,
					StartColumnIndex: column,
					EndColumnIndex:   column + 1,
					StartRowIndex:    fromRow,
				},
				Cell: &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						NumberFormat: &sheets.NumberFormat{
							Type:    "NUMBER",
							Pattern: "#,##0.00",
						},
					},
				},
				Fields: "userEnteredFormat",
			},
		},
	)

	valueRanges = append(valueRanges, &sheets.ValueRange{
		MajorDimension: "ROWS",
		Range:          fmt.Sprintf("'%s'!", title) + targetColumn + "4:" + fmt.Sprintf("%s%d", targetColumn, int64(len(sheet.RowData))),
		Values:         getColumnData(sheet, int(column), int(fromRow)),
	})

	return requests, valueRanges
}

var (
	baseCols1 = []string{"X", "AF", "AJ", "AP", "AT", "AX"}
	baseCols2 = []string{"U", "AA", "AG", "AM", "AR", "AV"}

	sums = map[string][]string{
		"Затраты на доставку транспортом":             append(baseCols1, "DU"),
		"Затраты на сертификацию предприятия":         append(baseCols2, "BG"),
		"Затраты на рекламу ИКУ за рубежом":           append(baseCols2, "BH"),
		"Затраты на перевод каталога ИКУ":             append(baseCols2, "BF"),
		"Затраты на аренду помещения ИКУ":             append(baseCols2, "BC"),
		"Затраты на сертификацию ИКУ":                 append(baseCols2, "BI"),
		"Затраты на демонстрацию ИКУ":                 append(baseCols2, "BC"),
		"Затраты на франчайзинг":                      append(baseCols2, "BC"),
		"Затраты на регистрацию товарных знаков":      append(baseCols2, "BO"),
		"Затраты на аренду":                           append(baseCols2, "BC"),
		"Затраты на перевод":                          append(baseCols2, "BF"),
		"Затраты на рекламу товаров за рубежом":       append(baseCols2, "BH"),
		"Затраты на участие в выставках":              append(baseCols2, "BV"),
		"Затраты на участие в выставках ИКУ":          append(baseCols2, "BV"),
		"Затраты на соответствие товаров требованиям": append(baseCols2, "BM"),
	}
)

func convertToFinancialSumCells_Expenses(batch *spreadsheetsAdapters.BatchUpdate, sheetID int64, sheet *sheets.GridData, title string) {
	args := sums[title]

	for _, arg := range args {
		requests, values := convertToFinancialSumCells(sheetID, arg, sheet, title)
		batch.WithRequest(requests...)
		batch.WithValueRange(values...)
	}
}
