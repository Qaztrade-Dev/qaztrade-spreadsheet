package adapters

import (
	"context"

	"google.golang.org/api/sheets/v4"
)

type BatchUpdate struct {
	service     *sheets.Service
	Requests    []*sheets.Request
	ValueRanges []*sheets.ValueRange
}

func NewBatchUpdate(service *sheets.Service) *BatchUpdate {
	return &BatchUpdate{
		service:     service,
		Requests:    make([]*sheets.Request, 0),
		ValueRanges: make([]*sheets.ValueRange, 0),
	}
}

func (b *BatchUpdate) WithProtectedRange(sheetID int64, protectedRanges []*sheets.ProtectedRange) {
	for _, pr := range protectedRanges {
		b.Requests = append(b.Requests, &sheets.Request{
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
	b.Requests = append(b.Requests, &sheets.Request{
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
	b.Requests = append(b.Requests, requests...)
}

func (b *BatchUpdate) WithValueRange(valueRanges ...*sheets.ValueRange) {
	b.ValueRanges = append(b.ValueRanges, valueRanges...)
}

func (b *BatchUpdate) Do(ctx context.Context, spreadsheetID string) error {
	if len(b.Requests) == 0 {
		return nil
	}
	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{Requests: b.Requests}
	if _, err := b.service.Spreadsheets.BatchUpdate(spreadsheetID, batchUpdateRequest).Context(ctx).Do(); err != nil {
		return err
	}

	if len(b.ValueRanges) > 0 {
		if _, err := b.service.Spreadsheets.Values.BatchUpdate(spreadsheetID, &sheets.BatchUpdateValuesRequest{
			Data:             b.ValueRanges,
			ValueInputOption: "USER_ENTERED",
		}).Context(ctx).Do(); err != nil {
			return err
		}
	}

	return nil
}
