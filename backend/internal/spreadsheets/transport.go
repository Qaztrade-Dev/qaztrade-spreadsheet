package spreadsheets

import (
	"net/http"

	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"
	authEndpoint "github.com/doodocs/qaztrade/backend/internal/auth/endpoint"
	authTransport "github.com/doodocs/qaztrade/backend/internal/auth/transport"

	spreadsheetsDomain "github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
	spreadsheetsEndpoint "github.com/doodocs/qaztrade/backend/internal/spreadsheets/endpoint"
	spreadsheetsService "github.com/doodocs/qaztrade/backend/internal/spreadsheets/service"
	spreadsheetsTransport "github.com/doodocs/qaztrade/backend/internal/spreadsheets/transport"

	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"

	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
)

func MakeHandler(svc spreadsheetsService.Service, jwtcli *jwt.Client, logger kitlog.Logger) http.Handler {
	var (
		opts = []kithttp.ServerOption{
			kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
			kithttp.ServerErrorEncoder(common.EncodeError),
			kithttp.ServerBefore(authTransport.WithRequestToken),
		}

		mdlwChainUser = endpoint.Chain(
			authEndpoint.MakeClaimsMiddleware[authDomain.UserClaims](jwtcli),
		)

		mdlwChainSpreadsheet = endpoint.Chain(
			authEndpoint.MakeClaimsMiddleware[spreadsheetsDomain.SpreadsheetClaims](jwtcli),
		)

		createSpreadsheetHandler = kithttp.NewServer(
			mdlwChainUser(
				spreadsheetsEndpoint.MakeCreateSpreadsheetEndpoint(svc),
			),
			spreadsheetsTransport.DecodeCreateSpreadsheetRequest, common.EncodeResponse,
			opts...,
		)
		listSpreadsheetsHandler = kithttp.NewServer(
			mdlwChainUser(
				spreadsheetsEndpoint.MakeListSpreadsheetsEndpoint(svc),
			),
			spreadsheetsTransport.DecodeListSpreadsheetsRequest, common.EncodeResponse,
			opts...,
		)
		addSheetHandler = kithttp.NewServer(
			mdlwChainSpreadsheet(
				spreadsheetsEndpoint.MakeAddSheetEndpoint(svc),
			),
			spreadsheetsTransport.DecodeAddSheetRequest, common.EncodeResponse,
			opts...,
		)
	)

	r := mux.NewRouter()
	r.Handle("/spreadsheets/", createSpreadsheetHandler).Methods("POST")
	r.Handle("/spreadsheets/", listSpreadsheetsHandler).Methods("GET")
	r.Handle("/spreadsheets/sheets/", addSheetHandler).Methods("POST")

	return r
}
