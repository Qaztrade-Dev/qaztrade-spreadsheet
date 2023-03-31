package spreadsheets

import (
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/doodocs/qaztrade/backend/internal/sheets/service"
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
		_ = opts

		// submitRecordHandler = kithttp.NewServer(
		// 	endpoint.MakeSubmitRecordEndpoint(svc, jwtcli),
		// 	sheetsTransport.DecodeSubmitRecordRequest, encodeResponse,
		// 	opts...,
		// )
		// submitApplicationHandler = kithttp.NewServer(
		// 	endpoint.MakeSubmitApplicationEndpoint(svc, jwtcli),
		// 	sheetsTransport.DecodeSubmitApplicationRequest, encodeResponse,
		// 	opts...,
		// )
		// addSheetHandler = kithttp.NewServer(
		// 	endpoint.MakeAddSheetEndpoint(svc, jwtcli),
		// 	sheetsTransport.DecodeAddSheetRequest, encodeResponse,
		// 	opts...,
		// )
	)

	r := mux.NewRouter()
	// r.Handle("/spreadsheets/", addSheetHandler).Methods("POST")
	/*
		POST /spreadsheets
		Authorization: Bearer <user_id:...>
	*/

	return r
}
