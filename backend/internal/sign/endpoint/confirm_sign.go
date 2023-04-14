package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sign/service"
	"github.com/go-kit/kit/endpoint"
)

type ConfirmSignRequest struct {
	DocumentID string
}

type ConfirmSignResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *ConfirmSignResponse) Error() error { return r.Err }

func MakeConfirmSignEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ConfirmSignRequest)

		err := s.ConfirmSign(ctx, &service.ConfirmSignRequest{
			SignDocumentID: req.DocumentID,
		})

		return ConfirmSignResponse{
			Err: err,
		}, nil
	}
}
