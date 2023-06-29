package manager

import (
	"net/http"

	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"
	authEndpoint "github.com/doodocs/qaztrade/backend/internal/auth/endpoint"
	authTransport "github.com/doodocs/qaztrade/backend/internal/auth/transport"
	"github.com/doodocs/qaztrade/backend/internal/common"
	managerEndpoint "github.com/doodocs/qaztrade/backend/internal/manager/endpoint"
	managerService "github.com/doodocs/qaztrade/backend/internal/manager/service"
	managerTransport "github.com/doodocs/qaztrade/backend/internal/manager/transport"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"

	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
)

func MakeHandler(svc managerService.Service, jwtcli *jwt.Client, logger kitlog.Logger) http.Handler {
	var (
		opts = []kithttp.ServerOption{
			kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
			kithttp.ServerErrorEncoder(common.EncodeError),
			kithttp.ServerBefore(authTransport.WithRequestToken),
		}

		mdlwChain = endpoint.Chain(
			authEndpoint.MakeClaimsMiddleware[authDomain.UserClaims](jwtcli),
			authEndpoint.MakeAuthMiddleware(authDomain.RoleManager),
		)

		switchStatusHandler = kithttp.NewServer(
			mdlwChain(managerEndpoint.MakeSwitchStatusEndpoint(svc)),
			managerTransport.DecodeSwitchStatusRequest, common.EncodeResponse,
			opts...,
		)

		listSpreadsheetsHandler = kithttp.NewServer(
			mdlwChain(managerEndpoint.MakeListSpreadsheetsEndpoint(svc)),
			managerTransport.DecodeListSpreadsheetsRequest, common.EncodeResponse,
			opts...,
		)

		downloadArchiveHandler = kithttp.NewServer(
			mdlwChain(managerEndpoint.MakeDownloadArchiveEndpoint(svc)),
			managerTransport.DecodeDownloadArchive, managerTransport.EncodeDownloadArchiveResponse,
			opts...,
		)

		getDDCardResponseHandler = kithttp.NewServer(
			mdlwChain(managerEndpoint.MakeGetDDCardResponseEndpoint(svc)),
			managerTransport.DecodeGetDDCardResponse, managerTransport.EncodeGetDDCardResponseResponse,
			opts...,
		)
	)

	r := mux.NewRouter()
	r.Handle("/manager/spreadsheets/", switchStatusHandler).Methods("PATCH")
	r.Handle("/manager/spreadsheets/", listSpreadsheetsHandler).Methods("GET")
	r.Handle("/manager/applications/{application_id}/archive", downloadArchiveHandler).Methods("GET")
	r.Handle("/manager/applications/{application_id}/ddcard", getDDCardResponseHandler).Methods("GET")

	return r
}
