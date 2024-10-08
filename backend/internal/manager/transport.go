package manager

import (
	"net/http"

	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"
	authEndpoint "github.com/doodocs/qaztrade/backend/internal/auth/endpoint"
	authTransport "github.com/doodocs/qaztrade/backend/internal/auth/transport"

	managerEndpoint "github.com/doodocs/qaztrade/backend/internal/manager/endpoint"
	managerService "github.com/doodocs/qaztrade/backend/internal/manager/service"
	managerTransport "github.com/doodocs/qaztrade/backend/internal/manager/transport"

	"github.com/doodocs/qaztrade/backend/internal/common"
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

		authChain = endpoint.Chain(
			authEndpoint.MakeClaimsMiddleware[authDomain.UserClaims](jwtcli),
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

		getDDCardHandler = kithttp.NewServer(
			mdlwChain(managerEndpoint.MakeGetDDCardEndpoint(svc)),
			managerTransport.DecodeGetDDCard, managerTransport.EncodeGetDDCardResponse,
			opts...,
		)
		getManagersHandler = kithttp.NewServer(
			authChain(managerEndpoint.MakeGetManagersEndpoint(svc)),
			managerTransport.DecodeGetManagers, common.EncodeResponse,
			opts...,
		)
		getNoticeHandler = kithttp.NewServer(
			mdlwChain(managerEndpoint.MakeGetNoticeEndpoint(svc)),
			managerTransport.DecodeGetNotice, managerTransport.EncodeGetNoticeResponse,
			opts...,
		)
		sendNoticeHandler = kithttp.NewServer(
			mdlwChain(managerEndpoint.MakeSendNoticeEndpoint(svc)),
			managerTransport.DecodeSendNotice, common.EncodeResponse,
			opts...,
		)
		grantPermissionsHandler = kithttp.NewServer(
			mdlwChain(managerEndpoint.MakeGrantPermissionsEndpoint(svc)),
			managerTransport.DecodeGrantPermissions, common.EncodeResponse,
			opts...,
		)
	)

	r := mux.NewRouter()
	r.Handle("/manager/spreadsheets/", switchStatusHandler).Methods("PATCH")
	r.Handle("/manager/spreadsheets/", listSpreadsheetsHandler).Methods("GET")
	r.Handle("/manager/applications/{application_id}/ddcard", getDDCardHandler).Methods("GET")
	r.Handle("/manager/applications/{application_id}/access", grantPermissionsHandler).Methods("POST")
	r.Handle("/manager/managers/", getManagersHandler).Methods("GET")
	r.Handle("/manager/{application_id}/notice", getNoticeHandler).Methods("GET")
	r.Handle("/manager/{application_id}/notice", sendNoticeHandler).Methods("POST")

	return r
}
