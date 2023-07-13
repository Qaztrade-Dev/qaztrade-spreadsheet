package adapters

import (
	"context"
	"fmt"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetClient struct {
	sheetsService *sheets.Service
}

var _ domain.SpreadsheetRepository = (*SpreadsheetClient)(nil)

func NewSpreadsheetClient(ctx context.Context, credentialsJson []byte) (*SpreadsheetClient, error) {
	sheetsService, err := sheets.NewService(
		ctx,
		option.WithCredentialsJSON(credentialsJson),
	)
	if err != nil {
		return nil, err
	}

	return &SpreadsheetClient{
		sheetsService: sheetsService,
	}, nil
}

func (c *SpreadsheetClient) GetSheetData(ctx context.Context, spreadsheetID string, sheetTitle string) ([][]string, error) {
	spreadsheet, err := c.sheetsService.Spreadsheets.Get(spreadsheetID).
		IncludeGridData(true).
		Ranges(sheetTitle).
		Context(ctx).
		Do()
	if err != nil {
		return nil, err
	}

	if len(spreadsheet.Sheets) == 0 {
		return nil, domain.ErrorSheetNotFound
	}

	var (
		sheet  = spreadsheet.Sheets[0]
		data   = sheet.Data[0].RowData
		result = make([][]string, len(data))
	)

	for i := range data {
		row := data[i]
		result[i] = make([]string, len(row.Values))

		for j := range row.Values {
			cell := row.Values[j]
			result[i][j] = getValue(cell.UserEnteredValue)
		}
	}

	return result, nil
}

func getValue(input *sheets.ExtendedValue) string {
	result := ""
	if input == nil {
		return result
	}

	switch {
	case input.BoolValue != nil:
		result = fmt.Sprintf("%v", *input.BoolValue)
	case input.NumberValue != nil:
		result = fmt.Sprintf("%f", *input.NumberValue)
	case input.StringValue != nil:
		result = *input.StringValue
	}

	return result
}
