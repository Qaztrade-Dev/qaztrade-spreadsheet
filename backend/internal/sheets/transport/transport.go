package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"github.com/doodocs/qaztrade/backend/internal/sheets/endpoint"
	"github.com/doodocs/qaztrade/backend/internal/sheets/pkg/tally"
)

func DecodeSubmitRecordRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		RowNumber int                    `json:"rowNum"`
		ParentID  string                 `json:"parentID"`
		ChildKey  string                 `json:"childKey"`
		Value     map[string]interface{} `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	var (
		tokenString = extractHeaderToken(r)
		sheetName   = extractHeaderSheetName(r)
	)

	sheetID, err := extractHeaderSheetID(r)
	if err != nil {
		return nil, err
	}

	return endpoint.SubmitRecordRequest{
		TokenString: tokenString,
		SheetName:   sheetName,
		SheetID:     sheetID,
		Payload: &domain.Payload{
			RowNumber: body.RowNumber,
			ParentID:  body.ParentID,
			ChildKey:  body.ChildKey,
			Value:     domain.PayloadValue(body.Value),
		},
	}, nil
}

func DecodeSubmitApplicationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	tallyJsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	application, err := tally.Decode(tallyJsonBytes)
	if err != nil {
		return nil, err
	}

	return endpoint.SubmitApplicationRequest{
		Application: application,
	}, nil
}

func DecodeAddSheetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		SheetName string `json:"sheet_name"`
	}
	fmt.Println("hello")

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	tokenString := extractHeaderToken(r)

	return endpoint.AddSheetRequest{
		TokenString: tokenString,
		SheetName:   body.SheetName,
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

func extractHeaderSheetID(r *http.Request) (int64, error) {
	sheetIDStr := r.Header.Get("x-sheet-id")
	sheetID, err := strconv.ParseInt(sheetIDStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return sheetID, nil
}

func extractHeaderSheetName(r *http.Request) string {
	sheetName := r.Header.Get("x-sheet-name")
	unescaped, _ := url.QueryUnescape(sheetName)
	return unescaped
}
