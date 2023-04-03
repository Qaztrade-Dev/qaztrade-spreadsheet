package google

import (
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/doodocs/qaztrade/backend/internal/google/endpoint"
	"github.com/doodocs/qaztrade/backend/internal/google/service"
	googleTransport "github.com/doodocs/qaztrade/backend/internal/google/transport"
	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"

	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
)

func MakeHandler(svc service.Service, logger kitlog.Logger) http.Handler {
	var (
		opts = []kithttp.ServerOption{
			kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
			kithttp.ServerErrorEncoder(common.EncodeError),
		}

		getRedirectLinkHandler = kithttp.NewServer(
			endpoint.MakeGetRedirectLinkEndpoint(svc),
			googleTransport.DecodeGetRedirectLinkRequest, googleTransport.EncodeGetRedirectLinkResponse,
			opts...,
		)

		updateTokenHandler = kithttp.NewServer(
			endpoint.MakeUpdateTokenEndpoint(svc),
			googleTransport.DecodeUpdateTokenRequest, common.EncodeResponse,
			opts...,
		)
	)

	r := mux.NewRouter()
	r.Handle("/google/auth", getRedirectLinkHandler).Methods(http.MethodGet)
	r.Handle("/google/callback", updateTokenHandler).Methods(http.MethodGet)

	return r
}
