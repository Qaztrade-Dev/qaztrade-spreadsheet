package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/service"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/go-kit/kit/endpoint"
)

type AddSheetRequest struct {
	TokenString string
	SheetName   string
}

type AddSheetResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *AddSheetResponse) Error() error { return r.Err }

func MakeAddSheetEndpoint(s service.Service, j *jwt.Client) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddSheetRequest)

		claims, err := jwt.Parse[domain.SpreadsheetClaims](j, req.TokenString)
		if err != nil {
			return nil, err
		}

		err = s.AddSheet(ctx, &service.AddSheetRequest{
			SpreadsheetID: claims.SpreadsheetID,
			SheetName:     req.SheetName,
		})

		return &AddSheetResponse{Err: err}, nil
	}
}
