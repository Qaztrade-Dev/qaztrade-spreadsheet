package adapters

import (
	"context"
	"fmt"

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

	deleteRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				DeleteDimension: &sheets.DeleteDimensionRequest{
					Range: &sheets.DimensionRange{
						SheetId:    1366326786,
						Dimension:  "ROWS",
						StartIndex: 200,
					},
				},
			},
		},
	}

	spreadsheetID := "1AfhwOrdMQMhmbgJNQhNOgt50fAmN1hl65WVFSmXqS68"
	_, err = spreadsheetsSvc.Spreadsheets.BatchUpdate(spreadsheetID, deleteRequest).Do()
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
