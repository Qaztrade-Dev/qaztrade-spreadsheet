package adapters

import (
	"context"

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

func (b *BatchUpdate) WithRequest(requests ...*sheets.Request) {
	b.requests = append(b.requests, requests...)
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

func (b *BatchUpdate) Do(ctx context.Context, spreadsheetID string) error {
	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{Requests: b.requests}
	_, err := b.service.Spreadsheets.BatchUpdate(spreadsheetID, batchUpdateRequest).Context(ctx).Do()
	return err
}
