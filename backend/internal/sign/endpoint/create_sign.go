package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sign/domain"
	"github.com/doodocs/qaztrade/backend/internal/sign/service"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/go-kit/kit/endpoint"
)

type CreateSignRequest struct {
	Token string
}

type CreateSignResponse struct {
	Link string `json:"link,omitempty"`
	Err  error  `json:"err,omitempty"`
}

func (r *CreateSignResponse) Error() error { return r.Err }

func MakeCreateSignEndpoint(s service.Service, j *jwt.Client) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateSignRequest)

		claims, err := jwt.Parse[domain.SpreadsheetClaims](j, req.Token)
		if err != nil {
			return nil, err
		}

		link, err := s.CreateSign(ctx, &service.CreateSignRequest{
			SpreadsheetID: claims.SpreadsheetID,
		})

		return &CreateSignResponse{
			Link: link,
			Err:  err,
		}, nil
	}
}
