package endpoint

import (
	"context"

	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/pkg/jsonspreadsheets"
	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/service"
	"github.com/go-kit/kit/endpoint"
)

type CreateSpreadsheetResponse struct {
	Link string `json:"link,omitempty"`
	Err  error  `json:"err,omitempty"`
}

func (r *CreateSpreadsheetResponse) Error() error { return r.Err }

func MakeCreateSpreadsheetEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		claims, err := authDomain.ExtractClaims[authDomain.UserClaims](ctx)
		if err != nil {
			return nil, err
		}

		link, err := s.CreateSpreadsheet(ctx, &service.CreateSpreadsheetRequest{
			UserID: claims.UserID,
		})

		return &CreateSpreadsheetResponse{
			Link: link,
			Err:  err,
		}, nil
	}
}

type ListSpreadsheetsRequest struct {
	Limit  uint64
	Offset uint64
}

type ListSpreadsheetsResponse struct {
	ApplicationList *jsonspreadsheets.ApplicationList `json:"list,omitempty"`
	Err             error                             `json:"err,omitempty"`
}

func (r *ListSpreadsheetsResponse) Error() error { return r.Err }

func MakeListSpreadsheetsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListSpreadsheetsRequest)

		claims, err := authDomain.ExtractClaims[authDomain.UserClaims](ctx)
		if err != nil {
			return nil, err
		}

		list, err := s.ListSpreadsheets(ctx, &service.ListSpreadsheetsRequest{
			UserID: claims.UserID,
			Limit:  req.Limit,
			Offset: req.Offset,
		})

		return &ListSpreadsheetsResponse{
			ApplicationList: jsonspreadsheets.EncodeApplicationList(list),
			Err:             err,
		}, nil
	}
}
