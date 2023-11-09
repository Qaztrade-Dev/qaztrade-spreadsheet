package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/assignments/pkg/jsondomain"
	"github.com/doodocs/qaztrade/backend/internal/assignments/service"
	"github.com/go-kit/kit/endpoint"
)

type GetAssignmentDialogResponse struct {
	Err      error                       `json:"err,omitempty"`
	Messages []*jsondomain.DialogMessage `json:"messages"`
}

func (r *GetAssignmentDialogResponse) Error() error { return r.Err }

func MakeGetAssignmentDialogEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		input := request.(uint64)

		response, err := s.GetAssignmentDialog(ctx, input)
		if err != nil {
			return &GetAssignmentDialogResponse{Err: err}, nil
		}

		return &GetAssignmentDialogResponse{
			Err:      err,
			Messages: jsondomain.EncodeSlice(response.Messages, jsondomain.EncodeDialogMessage),
		}, nil
	}
}
