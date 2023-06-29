package sign

import (
	"net/http"

	authEndpoint "github.com/doodocs/qaztrade/backend/internal/auth/endpoint"
	authTransport "github.com/doodocs/qaztrade/backend/internal/auth/transport"
	spreadsheetsDomain "github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"

	signEndpoint "github.com/doodocs/qaztrade/backend/internal/sign/endpoint"
	signService "github.com/doodocs/qaztrade/backend/internal/sign/service"
	signTransport "github.com/doodocs/qaztrade/backend/internal/sign/transport"

	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"

	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
)

func MakeHandler(svc signService.Service, jwtcli *jwt.Client, logger kitlog.Logger) http.Handler {
	var (
		opts = []kithttp.ServerOption{
			kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
			kithttp.ServerErrorEncoder(common.EncodeError),
			kithttp.ServerBefore(authTransport.WithRequestToken),
		}

		mdlwChainSpreadsheet = endpoint.Chain(
			authEndpoint.MakeClaimsMiddleware[spreadsheetsDomain.SpreadsheetClaims](jwtcli),
		)

		createSignHandler = kithttp.NewServer(
			mdlwChainSpreadsheet(
				signEndpoint.MakeCreateSignEndpoint(svc),
			),
			signTransport.DecodeCreateSignRequest, common.EncodeResponse,
			opts...,
		)

		confirmSignHandler = kithttp.NewServer(
			signEndpoint.MakeConfirmSignEndpoint(svc),
			signTransport.DecodeConfirmSignRequest, common.EncodeResponse,
			opts...,
		)

		syncSpreadsheetsHandler = kithttp.NewServer(
			signEndpoint.MakeSyncSpreadsheetsEndpoint(svc),
			signTransport.DecodeSyncSpreadsheetsRequest, common.EncodeResponse,
			opts...,
		)

		syncSigningTimeHandler = kithttp.NewServer(
			signEndpoint.MakeSyncSigningTimeEndpoint(svc),
			signTransport.DecodeSyncSigningTimeRequest, common.EncodeResponse,
			opts...,
		)
	)

	r := mux.NewRouter()
	r.Handle("/sign/", createSignHandler).Methods("POST")
	r.Handle("/sign/callback", confirmSignHandler).Methods("POST")
	r.Handle("/sign/sync/spreadsheets", syncSpreadsheetsHandler).Methods("POST")
	r.Handle("/sign/sync/sign", syncSigningTimeHandler).Methods("POST")

	return r
}
