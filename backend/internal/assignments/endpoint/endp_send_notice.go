package endpoint

import (
	"context"
	"io"

	"github.com/doodocs/qaztrade/backend/internal/assignments/service"
	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/go-kit/kit/endpoint"
)

type SendNoticeRequest struct {
	AssignmentID uint64
	FileReader   io.Reader
	FileSize     int64
	FileName     string
}

type SendNoticeResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *SendNoticeResponse) Error() error { return r.Err }

func MakeSendNoticeEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		input := request.(SendNoticeRequest)

		claims, err := authDomain.ExtractClaims[authDomain.UserClaims](ctx)
		if err != nil {
			return nil, err
		}

		err = s.SendNotice(ctx, &service.SendNoticeRequest{
			UserID:       claims.UserID,
			AssignmentID: input.AssignmentID,
			FileReader:   input.FileReader,
			FileSize:     input.FileSize,
			FileName:     input.FileName,
		})

		return &SendNoticeResponse{
			Err: err,
		}, nil
	}
}
