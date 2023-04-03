package spreadsheets

import (
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/endpoint"
	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/service"
	spreadsheetsTransport "github.com/doodocs/qaztrade/backend/internal/spreadsheets/transport"
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

		createSpreadsheetHandler = kithttp.NewServer(
			endpoint.MakeCreateSpreadsheetEndpoint(svc, jwtcli),
			spreadsheetsTransport.DecodeCreateSpreadsheetRequest, common.EncodeResponse,
			opts...,
		)

		listSpreadsheetsHandler = kithttp.NewServer(
			endpoint.MakeListSpreadsheetsEndpoint(svc, jwtcli),
			spreadsheetsTransport.DecodeListSpreadsheetsRequest, common.EncodeResponse,
			opts...,
		)
	)

	r := mux.NewRouter()
	r.Handle("/spreadsheets/", createSpreadsheetHandler).Methods("POST")
	r.Handle("/spreadsheets/", listSpreadsheetsHandler).Methods("GET")

	return r
}
