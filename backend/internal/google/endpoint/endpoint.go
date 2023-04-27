package endpoint

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/google/service"
	"github.com/go-kit/kit/endpoint"
)

type GetRedirectLinkResponse struct {
	Link string `json:"link,omitempty"`
	Err  error  `json:"err,omitempty"`
}

func (r *GetRedirectLinkResponse) Error() error { return r.Err }

func MakeGetRedirectLinkEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		redirectLink, err := s.GetRedirectLink(ctx)
		return &GetRedirectLinkResponse{Link: redirectLink, Err: err}, nil
	}
}

type UpdateTokenRequest struct {
	Code string
}

type UpdateTokenResponse struct {
	Err error `json:"err,omitempty"`
}

func (r *UpdateTokenResponse) Error() error { return r.Err }

func MakeUpdateTokenEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateTokenRequest)
		err := s.UpdateToken(ctx, req.Code)
		return &UpdateTokenResponse{Err: err}, nil
	}
}
