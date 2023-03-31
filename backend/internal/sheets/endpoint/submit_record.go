package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"github.com/doodocs/qaztrade/backend/internal/sheets/service"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/go-kit/kit/endpoint"
)

type SubmitRecordRequest struct {
	TokenString string
	SheetName   string
	SheetID     int64
	Payload     *domain.Payload
}

type SubmitRecordResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *SubmitRecordResponse) Error() error { return r.Err }

func MakeSubmitRecordEndpoint(s service.Service, j *jwt.Client) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SubmitRecordRequest)

		claims, err := j.Parse(req.TokenString)
		if err != nil {
			return nil, err
		}

		err = s.SubmitRecord(ctx, &service.SubmitRecordRequest{
			SpreadsheetID: claims.SpreadsheetID,
			SheetName:     req.SheetName,
			SheetID:       req.SheetID,
			Payload:       req.Payload,
		})
		return SubmitRecordResponse{Err: err}, nil
	}
}
