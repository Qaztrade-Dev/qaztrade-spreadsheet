package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"github.com/doodocs/qaztrade/backend/internal/sheets/service"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/go-kit/kit/endpoint"
)

type AddRowsRequest struct {
	Token string
	Input *domain.AddRowsInput
}

type AddRowsResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *AddRowsResponse) Error() error { return r.Err }

func MakeAddRowsEndpoint(s service.Service, j *jwt.Client) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRowsRequest)

		claims, err := jwt.Parse[domain.SpreadsheetClaims](j, req.Token)
		if err != nil {
			return nil, err
		}

		err = s.AddRows(ctx, &service.AddRowsRequest{
			SpreadsheetID: claims.SpreadsheetID,
			Input:         req.Input,
		})

		return AddRowsResponse{Err: err}, nil
	}
}
