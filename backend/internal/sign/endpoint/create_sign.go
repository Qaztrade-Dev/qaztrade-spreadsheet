package endpoint

import (
	"context"

	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/doodocs/qaztrade/backend/internal/sign/service"
	spreadheetsDomain "github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
	"github.com/go-kit/kit/endpoint"
)

type CreateSignResponse struct {
	Link string `json:"link,omitempty"`
	Err  error  `json:"err,omitempty"`
}

func (r *CreateSignResponse) Error() error { return r.Err }

func MakeCreateSignEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		claims, err := authDomain.ExtractClaims[spreadheetsDomain.SpreadsheetClaims](ctx)
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
