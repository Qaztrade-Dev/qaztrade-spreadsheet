package endpoint

import (
	"context"
	"io"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"github.com/doodocs/qaztrade/backend/internal/sheets/service"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/go-kit/kit/endpoint"
)

type UploadFileRequest struct {
	SpreadsheetTokenStr string
	SheetID             int64
	SheetName           string
	Hyperlink           string
	RowIdx              int64
	ColumnIdx           int64
	FileReader          io.Reader
	FileSize            int64
	FileName            string
}

type UploadFileResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *UploadFileResponse) Error() error { return r.Err }

func MakeUploadFileEndpoint(s service.Service, j *jwt.Client) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UploadFileRequest)

		claims, err := jwt.Parse[domain.SpreadsheetClaims](j, req.SpreadsheetTokenStr)
		if err != nil {
			return nil, err
		}

		err = s.UploadFile(ctx, &service.UploadFileRequest{
			SpreadsheetID: claims.SpreadsheetID,
			SheetID:       req.SheetID,
			SheetName:     req.SheetName,
			RowIdx:        req.RowIdx,
			ColumnIdx:     req.ColumnIdx,
			FileReader:    req.FileReader,
			FileSize:      req.FileSize,
			FileName:      req.FileName,
		})
		return &UploadFileResponse{Err: err}, nil
	}
}
