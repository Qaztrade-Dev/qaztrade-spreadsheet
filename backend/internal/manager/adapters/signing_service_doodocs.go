package adapters

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
)

type SigningServiceDoodocs struct {
	urlBase  string
	login    string
	password string
}

var _ domain.SigningService = (*SigningServiceDoodocs)(nil)

func NewSigningServiceDoodocs(urlBase, login, password string) *SigningServiceDoodocs {
	return &SigningServiceDoodocs{
		urlBase:  urlBase,
		login:    login,
		password: password,
	}
}

func (s *SigningServiceDoodocs) GetDDCard(ctx context.Context, documentID string) (*http.Response, error) {
	session, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return session.downloadDDCard(ctx, documentID)
}

type signingServiceDoodocsSession struct {
	urlBase string
	token   string
}

func (s *SigningServiceDoodocs) authenticate(ctx context.Context) (*signingServiceDoodocsSession, error) {
	var (
		url        = fmt.Sprintf("%s/api/v1/session", s.urlBase)
		payload    = strings.NewReader("grant_type=client_credentials")
		client     = &http.Client{}
		basicToken = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.login, s.password)))
	)

	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", basicToken))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var authResponse struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&authResponse); err != nil {
		return nil, err
	}

	result := &signingServiceDoodocsSession{
		urlBase: s.urlBase,
		token:   authResponse.AccessToken,
	}

	return result, nil
}

func (s *signingServiceDoodocsSession) downloadDDCard(ctx context.Context, documentID string) (*http.Response, error) {
	var (
		url = fmt.Sprintf("%s/api/v1/documents/%s/download", s.urlBase, documentID)
	)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
