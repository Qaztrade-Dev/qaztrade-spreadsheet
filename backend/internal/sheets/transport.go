package sheets

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"github.com/doodocs/qaztrade/backend/internal/sheets/service"
	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"

	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
)

func MakeHandler(svc service.Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	submitRecordHandler := kithttp.NewServer(
		makeSubmitRecordEndpoint(svc),
		decodeSubmitRecordRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/sheets/records", submitRecordHandler).Methods("POST")

	return r
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func decodeSubmitRecordRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		SpreadsheetID string `json:"spreadsheet_id"`
		Payload       struct {
			ParentID string                 `json:"parent_id"`
			ChildKey string                 `json:"child_key"`
			Value    map[string]interface{} `json:"value"`
		} `json:"payload"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return submitRecordRequest{
		SpreadsheetID: body.SpreadsheetID,
		Payload: &domain.Payload{
			ParentID: body.Payload.ParentID,
			ChildKey: body.Payload.ChildKey,
			Value:    domain.PayloadValue(body.Payload.Value),
		},
	}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
