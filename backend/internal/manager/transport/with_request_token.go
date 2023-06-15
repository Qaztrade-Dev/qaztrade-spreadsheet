package transport

import (
	"context"
	"net/http"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/manager/endpoint"
)

func WithRequestToken(ctx context.Context, req *http.Request) context.Context {
	var (
		tokenStr = extractHeaderToken(req)
		newCtx   = endpoint.WithToken(ctx, tokenStr)
	)

	return newCtx
}

func extractHeaderToken(r *http.Request) string {
	authorization := r.Header.Get("authorization")
	if authorization == "" {
		return ""
	}

	tokenString := strings.Split(authorization, " ")[1]
	return tokenString
}
