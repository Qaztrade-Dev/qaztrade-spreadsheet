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

func (b *BatchUpdate) Do(ctx context.Context, spreadsheetID string) error {
	if len(b.requests) == 0 {
		return nil
	}
	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{Requests: b.requests}
	_, err := b.service.Spreadsheets.BatchUpdate(spreadsheetID, batchUpdateRequest).Context(ctx).Do()
	return err
}
