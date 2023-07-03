package assignments

import (
	"net/http"

	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"
	authEndpoint "github.com/doodocs/qaztrade/backend/internal/auth/endpoint"
	authTransport "github.com/doodocs/qaztrade/backend/internal/auth/transport"

	assignmentsEndpoint "github.com/doodocs/qaztrade/backend/internal/assignments/endpoint"
	assignmentsService "github.com/doodocs/qaztrade/backend/internal/assignments/service"
	assignmentsTransport "github.com/doodocs/qaztrade/backend/internal/assignments/transport"

	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"

	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
)

func MakeHandler(svc assignmentsService.Service, jwtcli *jwt.Client, logger kitlog.Logger) http.Handler {
	var (
		opts = []kithttp.ServerOption{
			kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
			kithttp.ServerErrorEncoder(common.EncodeError),
			kithttp.ServerBefore(authTransport.WithRequestToken),
		}

		mdlwChainAdmin = endpoint.Chain(
			authEndpoint.MakeClaimsMiddleware[authDomain.UserClaims](jwtcli),
			authEndpoint.MakeAuthMiddleware(authDomain.RoleAdmin),
		)

		mdlwChainManager = endpoint.Chain(
			authEndpoint.MakeClaimsMiddleware[authDomain.UserClaims](jwtcli),
			authEndpoint.MakeAuthMiddleware(authDomain.RoleManager),
		)

		createBatchHandler = kithttp.NewServer(
			mdlwChainAdmin(assignmentsEndpoint.MakeCreateBatchEndpoint(svc)),
			assignmentsTransport.DecodeCreateBatchRequest, common.EncodeResponse,
			opts...,
		)

		getAssignmentsHandler = kithttp.NewServer(
			mdlwChainAdmin(assignmentsEndpoint.MakeGetAssignmentsEndpoint(svc)),
			assignmentsTransport.DecodeGetAssignmentsRequest, common.EncodeResponse,
			opts...,
		)

		getUserAssignmentsHandler = kithttp.NewServer(
			mdlwChainManager(assignmentsEndpoint.MakeGetUserAssignmentsEndpoint(svc)),
			assignmentsTransport.DecodeGetUserAssignmentsRequest, common.EncodeResponse,
			opts...,
		)
	)

	r := mux.NewRouter()
	r.Handle("/assignments/manager/", getUserAssignmentsHandler).Methods(http.MethodGet)
	r.Handle("/assignments/admin/", getAssignmentsHandler).Methods(http.MethodGet)
	r.Handle("/assignments/admin/batch/", createBatchHandler).Methods(http.MethodPost)

	return r
}
