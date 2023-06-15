package endpoint

import (
	"context"
	"io"
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
	"github.com/doodocs/qaztrade/backend/internal/manager/pkg/jsonmanager"
	"github.com/doodocs/qaztrade/backend/internal/manager/service"
	"github.com/go-kit/kit/endpoint"
)

type SwitchStatusRequest struct {
	ApplicationID string
	StatusName    string
}

type SwitchStatusResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *SwitchStatusResponse) Error() error { return r.Err }

func MakeSwitchStatusEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SwitchStatusRequest)

		err := s.SwitchStatus(ctx, &service.SwitchStatusRequest{
			ApplicationID: req.ApplicationID,
			StatusName:    req.StatusName,
		})

		return &SwitchStatusResponse{
			Err: err,
		}, nil
	}
}

type ListSpreadsheetsRequest struct {
	Limit  uint64
	Offset uint64
}

type ListSpreadsheetsResponse struct {
	ApplicationList *jsonmanager.ApplicationList `json:"list,omitempty"`
	Err             error                        `json:"err,omitempty"`
}

func (r *ListSpreadsheetsResponse) Error() error { return r.Err }

func MakeListSpreadsheetsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListSpreadsheetsRequest)

		list, err := s.ListSpreadsheets(ctx, &service.ListSpreadsheetsRequest{
			Limit:  req.Limit,
			Offset: req.Offset,
		})

		return &ListSpreadsheetsResponse{
			ApplicationList: jsonmanager.EncodeApplicationList(list),
			Err:             err,
		}, nil
	}
}

type DownloadArchiveRequest struct {
	ApplicationID string
}

type DownloadArchiveResponse struct {
	ArchiveReader io.ReadCloser
	RemoveFunc    domain.RemoveFunction
	Err           error `json:"err,omitempty"`
}

func (r *DownloadArchiveResponse) Error() error { return r.Err }

func MakeDownloadArchiveEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DownloadArchiveRequest)

		result, err := s.DownloadArchive(ctx, &service.DownloadArchiveRequest{
			ApplicationID: req.ApplicationID,
		})
		if err != nil {
			return nil, err
		}

		return &DownloadArchiveResponse{
			ArchiveReader: result.ArchiveReader,
			RemoveFunc:    result.RemoveFunc,
			Err:           err,
		}, nil
	}
}

type GetDDCardResponseRequest struct {
	ApplicationID string
}

type GetDDCardResponseResponse struct {
	HTTPResponse *http.Response
	Err          error `json:"err,omitempty"`
}

func (r *GetDDCardResponseResponse) Error() error { return r.Err }

func MakeGetDDCardResponseEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetDDCardResponseRequest)

		result, err := s.GetDDCardResponse(ctx, &service.GetDDCardResponseRequest{
			ApplicationID: req.ApplicationID,
		})
		if err != nil {
			return nil, err
		}

		return &GetDDCardResponseResponse{
			HTTPResponse: result,
			Err:          err,
		}, nil
	}
}
