package endpoint

import (
	"context"

	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/service"
	"github.com/go-kit/kit/endpoint"
)

type AddSheetRequest struct {
	SheetName string
}

type AddSheetResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *AddSheetResponse) Error() error { return r.Err }

func MakeAddSheetEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddSheetRequest)

		claims, err := authDomain.ExtractClaims[domain.SpreadsheetClaims](ctx)
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
