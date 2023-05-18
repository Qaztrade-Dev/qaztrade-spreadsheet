package endpoint

import (
	"context"
	"errors"
	"io"

	authDomain "github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
	"github.com/doodocs/qaztrade/backend/internal/manager/pkg/jsonmanager"
	"github.com/doodocs/qaztrade/backend/internal/manager/service"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/go-kit/kit/endpoint"
)

type SwitchStatusRequest struct {
	UserToken     string
	ApplicationID string
	StatusName    string
}

type SwitchStatusResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *SwitchStatusResponse) Error() error { return r.Err }

func MakeSwitchStatusEndpoint(s service.Service, j *jwt.Client) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SwitchStatusRequest)

		claims, err := jwt.Parse[authDomain.UserClaims](j, req.UserToken)
		if err != nil {
			return nil, err
		}

		if claims.Role != authDomain.RoleManager {
			return nil, errors.New("permission denied")
		}

		err = s.SwitchStatus(ctx, &service.SwitchStatusRequest{
			ApplicationID: req.ApplicationID,
			StatusName:    req.StatusName,
		})

		return &SwitchStatusResponse{
			Err: err,
		}, nil
	}
}

type ListSpreadsheetsRequest struct {
	UserToken string
	Limit     uint64
	Offset    uint64
}

type ListSpreadsheetsResponse struct {
	ApplicationList *jsonmanager.ApplicationList `json:"list,omitempty"`
	Err             error                        `json:"err,omitempty"`
}

func (r *ListSpreadsheetsResponse) Error() error { return r.Err }

func MakeListSpreadsheetsEndpoint(s service.Service, j *jwt.Client) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListSpreadsheetsRequest)

		claims, err := jwt.Parse[authDomain.UserClaims](j, req.UserToken)
		if err != nil {
			return nil, err
		}

		if claims.Role != authDomain.RoleManager {
			return nil, errors.New("permission denied")
		}

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
	UserToken     string
	ApplicationID string
}

type DownloadArchiveResponse struct {
	ArchiveReader io.ReadCloser
	RemoveFunc    domain.RemoveFunction
	Err           error `json:"err,omitempty"`
}

func (r *DownloadArchiveResponse) Error() error { return r.Err }

func MakeDownloadArchiveEndpoint(s service.Service, j *jwt.Client) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DownloadArchiveRequest)

		claims, err := jwt.Parse[authDomain.UserClaims](j, req.UserToken)
		if err != nil {
			return nil, err
		}

		if claims.Role != authDomain.RoleManager {
			return nil, errors.New("permission denied")
		}

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
