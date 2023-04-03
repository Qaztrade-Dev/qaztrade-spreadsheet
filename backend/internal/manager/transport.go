package manager

import (
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/doodocs/qaztrade/backend/internal/manager/endpoint"
	"github.com/doodocs/qaztrade/backend/internal/manager/service"
	managerTransport "github.com/doodocs/qaztrade/backend/internal/manager/transport"
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

		switchStatusHandler = kithttp.NewServer(
			endpoint.MakeSwitchStatusEndpoint(svc, jwtcli),
			managerTransport.DecodeSwitchStatusRequest, common.EncodeResponse,
			opts...,
		)

		listSpreadsheetsHandler = kithttp.NewServer(
			endpoint.MakeListSpreadsheetsEndpoint(svc, jwtcli),
			managerTransport.DecodeListSpreadsheetsRequest, common.EncodeResponse,
			opts...,
		)
	)

	r := mux.NewRouter()
	r.Handle("/manager/spreadsheets/", switchStatusHandler).Methods("PATCH")
	r.Handle("/manager/spreadsheets/", listSpreadsheetsHandler).Methods("GET")

	return r
}
