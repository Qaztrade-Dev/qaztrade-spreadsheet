package transport

import (
	"context"
	"net/http"

	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/doodocs/qaztrade/backend/internal/google/endpoint"
)

func DecodeGetRedirectLinkRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func EncodeGetRedirectLinkResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		common.EncodeError(ctx, e.Error(), w)
		return nil
	}
	resp := response.(*endpoint.GetRedirectLinkResponse)
	w.Header().Set("Location", resp.Link)
	w.WriteHeader(http.StatusTemporaryRedirect)
	return nil
}

func DecodeUpdateTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	code := r.URL.Query().Get("code")

	return endpoint.UpdateTokenRequest{
		Code: code,
	}, nil
}
