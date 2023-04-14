package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/sign/endpoint"
)

func DecodeCreateSignRequest(_ context.Context, r *http.Request) (interface{}, error) {
	tokenString := extractHeaderToken(r)

	return endpoint.CreateSignRequest{
		Token: tokenString,
	}, nil
}

func extractHeaderToken(r *http.Request) string {
	authorization := r.Header.Get("authorization")
	if authorization == "" {
		return ""
	}

	tokenString := strings.Split(authorization, " ")[1]
	return tokenString
}

func DecodeConfirmSignRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		DocumentID string `json:"document_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return endpoint.ConfirmSignRequest{
		DocumentID: body.DocumentID,
	}, nil
}
