package sheets

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/sheets/endpoint"
	"github.com/doodocs/qaztrade/backend/internal/sheets/pkg/jwt"
	"github.com/doodocs/qaztrade/backend/internal/sheets/service"
	sheetsTransport "github.com/doodocs/qaztrade/backend/internal/sheets/transport"
	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"

	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
)

func MakeHandler(svc service.Service, jwtcli *jwt.Client, logger kitlog.Logger) http.Handler {
	var (
		opts = []kithttp.ServerOption{
			kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
			kithttp.ServerErrorEncoder(encodeError),
		}

		submitRecordHandler = kithttp.NewServer(
			endpoint.MakeSubmitRecordEndpoint(svc, jwtcli),
			sheetsTransport.DecodeSubmitRecordRequest, encodeResponse,
			opts...,
		)
		submitApplicationHandler = kithttp.NewServer(
			endpoint.MakeSubmitApplicationEndpoint(svc, jwtcli),
			sheetsTransport.DecodeSubmitApplicationRequest, encodeResponse,
			opts...,
		)
		addSheetHandler = kithttp.NewServer(
			endpoint.MakeAddSheetEndpoint(svc, jwtcli),
			sheetsTransport.DecodeAddSheetRequest, encodeResponse,
			opts...,
		)
	)

	r := mux.NewRouter()
	r.Handle("/sheets/", addSheetHandler).Methods("POST")
	r.Handle("/sheets/records", submitRecordHandler).Methods("POST")
	r.Handle("/sheets/application", submitApplicationHandler).Methods("POST")

	return r
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.Error() != nil {
		encodeError(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	Error() error
}

// encodeError from business-logic
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
