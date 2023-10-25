package endpoint

import (
	"context"
	"io"

	"github.com/doodocs/qaztrade/backend/internal/assignments/service"
	"github.com/doodocs/qaztrade/backend/pkg/storage"
	"github.com/go-kit/kit/endpoint"
)

type GetArchiveRequest struct {
	AssignmentID uint64
}

type GetArchiveResponse struct {
	ArchiveReader io.ReadCloser
	RemoveFunc    storage.RemoveFunction
	Err           error `json:"err,omitempty"`
}

func (r *GetArchiveResponse) Error() error { return r.Err }

func MakeGetArchiveEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetArchiveRequest)

		result, err := s.GetArchive(ctx, &service.GetArchiveRequest{
			AssignmentID: req.AssignmentID,
		})
		if err != nil {
			return nil, err
		}

		return &GetArchiveResponse{
			ArchiveReader: result.ArchiveReader,
			RemoveFunc:    result.RemoveFunc,
			Err:           err,
		}, nil
	}
}
