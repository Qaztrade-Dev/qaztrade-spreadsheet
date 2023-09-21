package endpoint

import (
	"bytes"
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

type ListSpreadsheetsResponse struct {
	ApplicationList *jsonmanager.ApplicationList `json:"list,omitempty"`
	Err             error                        `json:"err,omitempty"`
}

func (r *ListSpreadsheetsResponse) Error() error { return r.Err }

func MakeListSpreadsheetsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(domain.GetManyInput)

		list, err := s.ListSpreadsheets(ctx, &req)

		return &ListSpreadsheetsResponse{
			ApplicationList: jsonmanager.EncodeApplicationList(list),
			Err:             err,
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

type GetNoticeRequest struct {
	ApplicationID string
}

type GetNoticeResponse struct {
	Docx *bytes.Buffer
	Err  error `json:"err,omitempty"`
}

type SendNoticeRequest struct {
	ApplicationID string
	FileReader    io.Reader
	FileSize      int64
	FileName      string
}

type SendNoticeResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *GetDDCardResponse) Error() error { return r.Err }

func (r *GetNoticeResponse) Error() error { return r.Err }

func (r *SendNoticeResponse) Error() error { return r.Err }

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

func MakeGetNoticeEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetNoticeRequest)
		data, err := s.GetNotice(ctx, &service.GetNoticeRequest{
			ApplicationID: req.ApplicationID,
		})

		if err != nil {
			return nil, err
		}

		return &GetNoticeResponse{
			Docx: data,
			Err:  err,
		}, err
	}
}

func MakeSendNoticeEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SendNoticeRequest)
		err := s.SendNotice(ctx, &service.SendNoticeRequest{
			ApplicationID: req.ApplicationID,
			FileReader:    req.FileReader,
			FileSize:      req.FileSize,
			FileName:      req.FileName,
		})
		return &SendNoticeResponse{
			Err: err,
		}, err
	}
}

type GrantPermissionsRequest struct {
	ApplicationID string
}

type GrantPermissionsResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *GrantPermissionsResponse) Error() error { return r.Err }

func MakeGrantPermissionsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GrantPermissionsRequest)

		err := s.GrantPermissions(ctx, &service.GrantPermissionsRequest{
			ApplicationID: req.ApplicationID,
		})

		return &GrantPermissionsResponse{
			Err: err,
		}, nil
	}
}
