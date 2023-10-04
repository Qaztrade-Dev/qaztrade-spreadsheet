package adapters

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func (s *SpreadsheetServiceGoogle) FixEmptyExpensesSum(ctx context.Context) error {
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

	// spreadsheetIDs := []string{
	// 	"1DMmoKcY0PKck4W-4Im4HUj2ba7k4spd1FV2zCHhqSRw",
	// }

	for i := 935; i < len(spreadsheetIDs); i++ {
		spreadsheetID := spreadsheetIDs[i]
		fmt.Println(spreadsheetID)

		data, err := getDataFromRanges(ctx, spreadsheetsSvc, spreadsheetID, []string{
			"'Затраты на доставку транспортом'!X2",
			"'Затраты на доставку транспортом'!DU2",
		})
		if err != nil {
			if strings.Contains(err.Error(), "Unable to parse range") {
				continue
			}

			return err
		}

		var (
			avrSum      = extractValue(data[0])
			expensesSum = extractValue(data[1])
		)

		if avrSum > 0 && expensesSum == 0 {
			fmt.Printf("FOUND: %s\n", spreadsheetID)
		}

		fmt.Printf("%v/%v\n", i+1, len(spreadsheetIDs))
		time.Sleep(time.Second)
	}

	return nil
}

func extractValue(batchDataValue [][]interface{}) float64 {
	var value string
	if len(batchDataValue) > 0 && len(batchDataValue[0]) > 0 {
		value = strings.TrimSpace(batchDataValue[0][0].(string))
	}

	value = strings.ReplaceAll(value, "\u00a0", "")
	value = strings.ReplaceAll(value, ",", ".")
	value = strings.ReplaceAll(value, " ", "")

	floatValue, _ := strconv.ParseFloat(value, 64)
	return floatValue
}

func getDataFromRanges(ctx context.Context, sheetsSvc *sheets.Service, spreadsheetID string, ranges []string) ([][][]interface{}, error) {
	resp, err := sheetsSvc.Spreadsheets.Values.BatchGet(spreadsheetID).Ranges(ranges...).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	datas := make([][][]interface{}, len(resp.ValueRanges))
	for i := range resp.ValueRanges {
		datas[i] = resp.ValueRanges[i].Values
	}

	return datas, nil
}
