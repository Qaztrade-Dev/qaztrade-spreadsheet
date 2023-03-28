package transport

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

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
	spreadsheetID := "1KL-lrhs-Wu9kRAppBxAHUUFr7OCfNYla8Z7W-0tX4Mo"

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

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

	// TODO
	// parse spreadsheetID from jwt
	spreadsheetID := "1I7tYAhUjPJGaMU7_XbhC08rQw55IRc7bEtg1mgmRPKg"

	return endpoint.SubmitApplicationRequest{
		SpreadsheetID: spreadsheetID,
		Application:   application,
	}, nil
}
