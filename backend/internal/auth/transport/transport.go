package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/auth/endpoint"
)

func DecodeSignUpRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		OrgName  string `json:"org_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return endpoint.SignUpRequest{
		Email:    body.Email,
		Password: body.Password,
		OrgName:  body.OrgName,
	}, nil
}

func DecodeSignInRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return endpoint.SignInRequest{
		Email:    body.Email,
		Password: body.Password,
	}, nil
}

func DecodeForgotRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return endpoint.ForgotRequest{
		Email: body.Email,
	}, nil
}

func DecodeRestoreRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	tokenString := extractHeaderToken(r)

	return endpoint.RestoreRequest{
		AccessToken: tokenString,
		Password:    body.Password,
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