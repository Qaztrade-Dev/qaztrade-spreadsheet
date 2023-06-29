package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/endpoint"
)

func DecodeCreateSpreadsheetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func DecodeListSpreadsheetsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		limitStr  = r.URL.Query().Get("limit")
		offsetStr = r.URL.Query().Get("offset")

		limit, _  = strconv.ParseUint(limitStr, 10, 0)
		offset, _ = strconv.ParseUint(offsetStr, 10, 0)
	)

	return endpoint.ListSpreadsheetsRequest{
		Limit:  limit,
		Offset: offset,
	}, nil
}

func DecodeAddSheetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		SheetName string `json:"sheet_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return endpoint.AddSheetRequest{
		SheetName: body.SheetName,
	}, nil
}
