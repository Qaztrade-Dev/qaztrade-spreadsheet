package endpoint

import (
	"context"
	"io"

	"github.com/doodocs/qaztrade/backend/internal/assignments/service"
	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/go-kit/kit/endpoint"
)

type RespondNoticeRequest struct {
	ApplicationID  string
	AssignmentType string
	FileReader     io.Reader
	FileSize       int64
	FileName       string
}

type RespondNoticeResponse struct {
	Err  error  `json:"err,omitempty"`
	Link string `json:"link,omitempty"`
}

func (r *RespondNoticeResponse) Error() error { return r.Err }

func MakeRespondNoticeEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		input := request.(RespondNoticeRequest)

		claims, err := authDomain.ExtractClaims[authDomain.UserClaims](ctx)
		if err != nil {
			return nil, err
		}

		response, err := s.RespondNotice(ctx, &service.RespondNoticeRequest{
			UserID:         claims.UserID,
			ApplicationID:  input.ApplicationID,
			AssignmentType: input.AssignmentType,
			FileReader:     input.FileReader,
			FileSize:       input.FileSize,
			FileName:       input.FileName,
		})

		return &RespondNoticeResponse{
			Err:  err,
			Link: response.SignLink,
		}, nil
	}
}
