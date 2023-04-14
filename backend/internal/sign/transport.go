package sign

import (
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/doodocs/qaztrade/backend/internal/sign/endpoint"
	"github.com/doodocs/qaztrade/backend/internal/sign/service"
	signTransport "github.com/doodocs/qaztrade/backend/internal/sign/transport"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"

	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
)

func MakeHandler(svc service.Service, jwtcli *jwt.Client, logger kitlog.Logger) http.Handler {
	var (
		opts = []kithttp.ServerOption{
			kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
			kithttp.ServerErrorEncoder(common.EncodeError),
		}

		createSignHandler = kithttp.NewServer(
			endpoint.MakeCreateSignEndpoint(svc, jwtcli),
			signTransport.DecodeCreateSpreadsheetRequest, common.EncodeResponse,
			opts...,
		)
	)

	r := mux.NewRouter()
	r.Handle("/sign/", createSignHandler).Methods("POST")

	return r
}
