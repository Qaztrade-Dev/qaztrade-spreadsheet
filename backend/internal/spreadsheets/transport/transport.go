package transport

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/endpoint"
)

func DecodeCreateSpreadsheetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	tokenString := extractHeaderToken(r)

	return endpoint.CreateSpreadsheetRequest{
		UserToken: tokenString,
	}, nil
}

func extractHeaderToken(r *http.Request) string {
	authorization := r.Header.Get("authorization")
	if authorization == "" {
		return ""
	}

	tokenString := strings.Split(authorization, " ")[1]
	return tokenString
}

func DecodeListSpreadsheetsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		limitStr    = r.URL.Query().Get("limit")
		offsetStr   = r.URL.Query().Get("offset")
		tokenString = extractHeaderToken(r)

		limit, _  = strconv.ParseUint(limitStr, 10, 0)
		offset, _ = strconv.ParseUint(offsetStr, 10, 0)
	)

	return endpoint.ListSpreadsheetsRequest{
		UserToken: tokenString,
		Limit:     limit,
		Offset:    offset,
	}, nil
}
