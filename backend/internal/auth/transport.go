package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/auth/endpoint"
	"github.com/doodocs/qaztrade/backend/internal/auth/service"
	authTransport "github.com/doodocs/qaztrade/backend/internal/auth/transport"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/gorilla/mux"

	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
)

func MakeHandler(svc service.Service, jwtcli *jwt.Client, logger kitlog.Logger) http.Handler {
	var (
		opts = []kithttp.ServerOption{
			kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
			kithttp.ServerErrorEncoder(encodeError),
		}

		signUpHandler = kithttp.NewServer(
			endpoint.MakeSignUpEndpoint(svc),
			authTransport.DecodeSignUpRequest, encodeResponse,
			opts...,
		)
		signInHandler = kithttp.NewServer(
			endpoint.MakeSignInEndpoint(svc),
			authTransport.DecodeSignInRequest, encodeResponse,
			opts...,
		)
		forgotHandler = kithttp.NewServer(
			endpoint.MakeForgotEndpoint(svc),
			authTransport.DecodeForgotRequest, encodeResponse,
			opts...,
		)
		restoreHandler = kithttp.NewServer(
			endpoint.MakeRestoreEndpoint(svc, jwtcli),
			authTransport.DecodeRestoreRequest, encodeResponse,
			opts...,
		)
	)

	r := mux.NewRouter()
	r.Handle("/auth/signup", signUpHandler).Methods("POST")
	r.Handle("/auth/signin", signInHandler).Methods("POST")
	r.Handle("/auth/forgot", forgotHandler).Methods("POST")
	r.Handle("/auth/restore", restoreHandler).Methods("POST")

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
