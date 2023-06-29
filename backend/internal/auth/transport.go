package auth

import (
	"net/http"

	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"
	authEndpoint "github.com/doodocs/qaztrade/backend/internal/auth/endpoint"
	authService "github.com/doodocs/qaztrade/backend/internal/auth/service"
	authTransport "github.com/doodocs/qaztrade/backend/internal/auth/transport"
	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/gorilla/mux"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
)

func MakeHandler(svc authService.Service, jwtcli *jwt.Client, logger kitlog.Logger) http.Handler {
	var (
		opts = []kithttp.ServerOption{
			kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
			kithttp.ServerErrorEncoder(common.EncodeError),
			kithttp.ServerBefore(authTransport.WithRequestToken),
		}

		mdlwChain = endpoint.Chain(
			authEndpoint.MakeClaimsMiddleware[authDomain.UserClaims](jwtcli),
		)

		signUpHandler = kithttp.NewServer(
			authEndpoint.MakeSignUpEndpoint(svc),
			authTransport.DecodeSignUpRequest, common.EncodeResponse,
			opts...,
		)
		signInHandler = kithttp.NewServer(
			authEndpoint.MakeSignInEndpoint(svc),
			authTransport.DecodeSignInRequest, common.EncodeResponse,
			opts...,
		)
		forgotHandler = kithttp.NewServer(
			authEndpoint.MakeForgotEndpoint(svc),
			authTransport.DecodeForgotRequest, common.EncodeResponse,
			opts...,
		)
		restoreHandler = kithttp.NewServer(
			mdlwChain(authEndpoint.MakeRestoreEndpoint(svc)),
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
