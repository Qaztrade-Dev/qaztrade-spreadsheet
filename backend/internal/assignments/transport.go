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

		mdlwChainUser = endpoint.Chain(
			authEndpoint.MakeClaimsMiddleware[authDomain.UserClaims](jwtcli),
			authEndpoint.MakeAuthMiddleware(authDomain.RoleUser),
		)

		createBatchHandler = kithttp.NewServer(
			mdlwChainAdmin(assignmentsEndpoint.MakeCreateBatchEndpoint(svc)),
			assignmentsTransport.DecodeCreateBatchRequest, common.EncodeResponse,
			opts...,
		)

		redistributeAssignmentsHandler = kithttp.NewServer(
			mdlwChainAdmin(
				assignmentsEndpoint.MakeRedistributeAssignmentsEndpoint(svc),
			),
			assignmentsTransport.DecodeRedistributeAssignmentsRequest, common.EncodeResponse,
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

		changeAssigneeHandler = kithttp.NewServer(
			mdlwChainAdmin(assignmentsEndpoint.MakeChangeAssigneeEndpoint(svc)),
			assignmentsTransport.DecodeChangeAssigneeRequest, common.EncodeResponse,
			opts...,
		)

		getArchiveHandler = kithttp.NewServer(
			mdlwChainManager(assignmentsEndpoint.MakeGetArchiveEndpoint(svc)),
			assignmentsTransport.DecodeGetArchiveRequest, assignmentsTransport.EncodeGetArchiveResponse,
			opts...,
		)

		checkAssignmentHandler = kithttp.NewServer(
			mdlwChainManager(assignmentsEndpoint.MakeCheckAssignmentEndpoint(svc)),
			assignmentsTransport.DecodeCheckAssignmentRequest, common.EncodeResponse,
			opts...,
		)

		enqueueAssignmentsHandler = kithttp.NewServer(
			assignmentsEndpoint.MakeEnqueueAssignmentsEndpoint(svc),
			assignmentsTransport.DecodeEnqueueAssignmentsRequest, common.EncodeResponse,
			opts...,
		)

		sendNoticeHandler = kithttp.NewServer(
			mdlwChainManager(assignmentsEndpoint.MakeSendNoticeEndpoint(svc)),
			assignmentsTransport.DecodeSendNotice, common.EncodeResponse,
			opts...,
		)

		respondNoticeHandler = kithttp.NewServer(
			mdlwChainUser(assignmentsEndpoint.MakeRespondNoticeEndpoint(svc)),
			assignmentsTransport.DecodeRespondNotice, common.EncodeResponse,
			opts...,
		)

		updateAssignmentStatusHandler = kithttp.NewServer(
			mdlwChainManager(assignmentsEndpoint.MakeUpdateAssignmentStatusEndpoint(svc)),
			assignmentsTransport.DecodeUpdateAssignmentStatusRequest, common.EncodeResponse,
			opts...,
		)

		getAssignmentDialogHandler = kithttp.NewServer(
			mdlwChainManager(assignmentsEndpoint.MakeGetAssignmentDialogEndpoint(svc)),
			assignmentsTransport.DecodeGetAssignmentDialogRequest, common.EncodeResponse,
			opts...,
		)
	)

	r := mux.NewRouter()
	r.Handle("/assignments/manager/", getUserAssignmentsHandler).Methods(http.MethodGet)
	r.Handle("/assignments/admin/", getAssignmentsHandler).Methods(http.MethodGet)
	r.Handle("/assignments/admin/batch/", createBatchHandler).Methods(http.MethodPost)
	r.Handle("/assignments/admin/batch/redistribute", redistributeAssignmentsHandler).Methods(http.MethodPost)
	r.Handle("/assignments/{assignment_id}/user", changeAssigneeHandler).Methods(http.MethodPatch)
	r.Handle("/assignments/{assignment_id}/archive", getArchiveHandler).Methods(http.MethodGet)
	r.Handle("/assignments/{assignment_id}/check", checkAssignmentHandler).Methods(http.MethodPost)
	r.Handle("/assignments/enqueue", enqueueAssignmentsHandler).Methods(http.MethodPost)
	r.Handle("/assignments/{assignment_id}/notice", sendNoticeHandler).Methods(http.MethodPost)
	r.Handle("/assignments/applications/{application_id}/replies/{assignment_type}", respondNoticeHandler).Methods(http.MethodPost)
	r.Handle("/assignments/{assignment_id}/status", updateAssignmentStatusHandler).Methods(http.MethodPatch)
	r.Handle("/assignments/{assignment_id}/dialog", getAssignmentDialogHandler).Methods(http.MethodGet)

	return r
}
