package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/manager/endpoint"
)

func DecodeSwitchStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		ApplicationID string `json:"application_id"`
		StatusName    string `json:"status_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	tokenString := extractHeaderToken(r)

	return endpoint.SwitchStatusRequest{
		UserToken:     tokenString,
		ApplicationID: body.ApplicationID,
		StatusName:    body.StatusName,
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
