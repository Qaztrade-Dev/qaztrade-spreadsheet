package sheets

import (
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/doodocs/qaztrade/backend/internal/sheets/endpoint"
	"github.com/doodocs/qaztrade/backend/internal/sheets/service"
	sheetsTransport "github.com/doodocs/qaztrade/backend/internal/sheets/transport"
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

		submitApplicationHandler = kithttp.NewServer(
			endpoint.MakeSubmitApplicationEndpoint(svc, jwtcli),
			sheetsTransport.DecodeSubmitApplicationRequest, common.EncodeResponse,
			opts...,
		)

		uploadFileHandler = kithttp.NewServer(
			endpoint.MakeUploadFileEndpoint(svc, jwtcli),
			sheetsTransport.DecodeUploadFileRequest, common.EncodeResponse,
			opts...,
		)
	)

	r := mux.NewRouter()
	r.Handle("/sheets/application", submitApplicationHandler).Methods("POST")
	r.Handle("/sheets/file", uploadFileHandler).Methods("POST")

	return r
}
