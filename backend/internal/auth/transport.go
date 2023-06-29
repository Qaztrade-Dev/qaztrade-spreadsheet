package auth

import (
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/auth/endpoint"
	"github.com/doodocs/qaztrade/backend/internal/auth/service"
	authTransport "github.com/doodocs/qaztrade/backend/internal/auth/transport"
	"github.com/doodocs/qaztrade/backend/internal/common"
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
			kithttp.ServerErrorEncoder(common.EncodeError),
			kithttp.ServerBefore(authTransport.WithRequestToken),
		}

		signUpHandler = kithttp.NewServer(
			endpoint.MakeSignUpEndpoint(svc),
			authTransport.DecodeSignUpRequest, common.EncodeResponse,
			opts...,
		)
		signInHandler = kithttp.NewServer(
			endpoint.MakeSignInEndpoint(svc),
			authTransport.DecodeSignInRequest, common.EncodeResponse,
			opts...,
		)
		forgotHandler = kithttp.NewServer(
			endpoint.MakeForgotEndpoint(svc),
			authTransport.DecodeForgotRequest, common.EncodeResponse,
			opts...,
		)
		restoreHandler = kithttp.NewServer(
			endpoint.MakeRestoreEndpoint(svc, jwtcli),
			authTransport.DecodeRestoreRequest, common.EncodeResponse,
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
