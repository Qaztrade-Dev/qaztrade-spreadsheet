package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/auth/service"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/log"
)

func MakeHandler(svc service.Service, jwtcli *jwt.Client, logger kitlog.Logger) http.Handler {
	var (
	// opts = []kithttp.ServerOption{
	// 	kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	// 	kithttp.ServerErrorEncoder(encodeError),
	// }

	// submitRecordHandler = kithttp.NewServer(
	// 	endpoint.MakeSubmitRecordEndpoint(svc, jwtcli),
	// 	sheetsTransport.DecodeSubmitRecordRequest, encodeResponse,
	// 	opts...,
	// )
	// submitApplicationHandler = kithttp.NewServer(
	// 	endpoint.MakeSubmitApplicationEndpoint(svc, jwtcli),
	// 	sheetsTransport.DecodeSubmitApplicationRequest, encodeResponse,
	// 	opts...,
	// )
	// addSheetHandler = kithttp.NewServer(
	// 	endpoint.MakeAddSheetEndpoint(svc, jwtcli),
	// 	sheetsTransport.DecodeAddSheetRequest, encodeResponse,
	// 	opts...,
	// )
	)

	r := mux.NewRouter()
	/*
		POST /auth/signup
		{
			"email": "example@example.com",
			"password": "password123",
			"org_name": "OpenAI Inc."
		}
		Resp:
		{"access_token": "..."}

		POST /auth/signin
		{
			"email": "example@example.com",
			"password": "password123"
		}
		Resp:
		{"access_token": "..."}

		POST /auth/forgot
		{
			"email": "example@example.com"
		}
		Resp:
		200 OK

		POST /auth/restore
		"Authorization": "Bearer ..."

		{
			"password": "..."
		}
		Resp:
		200 OK - Login using your password
	*/

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
