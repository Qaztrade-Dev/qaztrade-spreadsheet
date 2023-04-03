package transport

import (
	"context"
	"net/http"
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
