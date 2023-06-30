package endpoint

import (
	"context"
	"io"
	"net/http"
	"time"

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
	Limit            uint64
	Offset           uint64
	BIN              string
	CompensationType string
	SignedAtFrom     time.Time
	SignedAtUntil    time.Time
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
			Limit:            req.Limit,
			Offset:           req.Offset,
			BIN:              req.BIN,
			CompensationType: req.CompensationType,
			SignedAtFrom:     req.SignedAtFrom,
			SignedAtUntil:    req.SignedAtUntil,
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

type GetDDCardRequest struct {
	ApplicationID string
}

type GetDDCardResponse struct {
	HTTPResponse *http.Response
	Err          error `json:"err,omitempty"`
}

func (r *GetDDCardResponse) Error() error { return r.Err }

func MakeGetDDCardEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetDDCardRequest)

		result, err := s.GetDDCard(ctx, &service.GetDDCardRequest{
			ApplicationID: req.ApplicationID,
		})
		if err != nil {
			return nil, err
		}

		return &GetDDCardResponse{
			HTTPResponse: result,
			Err:          err,
		}, nil
	}
}

type GetManagersResponse struct {
	Managers []*jsonmanager.Manager `json:"managers"`
	Err      error                  `json:"err,omitempty"`
}

func (r *GetManagersResponse) Error() error { return r.Err }

func MakeGetManagersEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		result, err := s.GetManagers(ctx)
		if err != nil {
			return nil, err
		}

		return &GetManagersResponse{
			Managers: jsonmanager.EncodeSlice(result, jsonmanager.EncodeManager),
			Err:      err,
		}, nil
	}
}
