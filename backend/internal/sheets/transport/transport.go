package transport

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"github.com/doodocs/qaztrade/backend/internal/sheets/endpoint"
	"github.com/doodocs/qaztrade/backend/internal/sheets/pkg/tally"
)

func DecodeSubmitRecordRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		ParentID string                 `json:"parentID"`
		ChildKey string                 `json:"childKey"`
		Value    map[string]interface{} `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	spreadsheetID := "1KL-lrhs-Wu9kRAppBxAHUUFr7OCfNYla8Z7W-0tX4Mo"

	return endpoint.SubmitRecordRequest{
		SpreadsheetID: spreadsheetID,
		Payload: &domain.Payload{
			ParentID: body.ParentID,
			ChildKey: body.ChildKey,
			Value:    domain.PayloadValue(body.Value),
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

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	tokenString := extractToken(r)

	return endpoint.AddSheetRequest{
		TokenString: tokenString,
		SheetName:   body.SheetName,
	}, nil
}

func extractToken(r *http.Request) string {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return ""
	}

	tokenString := strings.Split(authorization, " ")[1]
	return tokenString
}
