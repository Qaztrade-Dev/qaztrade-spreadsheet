package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/sign/endpoint"
)

func DecodeCreateSignRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func DecodeConfirmSignRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		DocumentID string `json:"document_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return endpoint.ConfirmSignRequest{
		DocumentID: body.DocumentID,
	}, nil
}

func DecodeSyncSigningTimeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		DocumentID string `json:"document_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return endpoint.SyncSigningTimeRequest{
		DocumentID: body.DocumentID,
	}, nil
}

func DecodeSyncSpreadsheetsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		SpreadsheetID string `json:"spreadsheet_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return endpoint.SyncSpreadsheetsRequest{
		SpreadsheetID: body.SpreadsheetID,
	}, nil
}
