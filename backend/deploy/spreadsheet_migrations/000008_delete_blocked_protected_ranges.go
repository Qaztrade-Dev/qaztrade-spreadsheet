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

func (s *SpreadsheetServiceGoogle) DeleteBlockedProtectedRanges(ctx context.Context) error {
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

	for i := 1016; i < len(spreadsheetIDs); i++ {
		spreadsheetID := spreadsheetIDs[i]
		fmt.Println(spreadsheetID)

		spreadsheet, err := spreadsheetsSvc.Spreadsheets.Get(spreadsheetID).Context(ctx).Do()
		if err != nil {
			return err
		}

		batch := spreadsheetsAdapters.NewBatchUpdate(spreadsheetsSvc)

		for _, sheet := range spreadsheet.Sheets {
			for _, protectedRange := range sheet.ProtectedRanges {
				if strings.Contains(protectedRange.Description, "_blocked_") {
					batch.WithRequest(
						&sheets.Request{
							DeleteProtectedRange: &sheets.DeleteProtectedRangeRequest{
								ProtectedRangeId: protectedRange.ProtectedRangeId,
							},
						},
					)
				}
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
