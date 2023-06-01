package adapters

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/sign/domain"
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

func (s *SigningServiceDoodocs) CreateSigningDocument(ctx context.Context, documentName string, documentReader io.Reader) (*domain.CreateSigningDocumentResponse, error) {
	session, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	documentID, err := session.uploadPDF(ctx, documentName, documentReader)
	if err != nil {
		return nil, err
	}

	if err := session.configureWorkflow(ctx, documentID); err != nil {
		return nil, err
	}

	if err := session.launchDocument(ctx, documentID); err != nil {
		return nil, err
	}

	recipients, err := session.getRecipients(ctx, documentID)
	if err != nil {
		return nil, err
	}

	link, err := session.getRecipientLink(ctx, documentID, recipients.Recipients[0].OriginID)
	if err != nil {
		return nil, err
	}

	result := &domain.CreateSigningDocumentResponse{
		DocumentID: documentID,
		SignLink:   link,
	}

	return result, nil
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

func (s *signingServiceDoodocsSession) uploadPDF(ctx context.Context, documentName string, fileReader io.Reader) (string, error) {
	var (
		teamspaceID = "09852313-811c-43aa-bd63-529b3cf539af"
		url         = fmt.Sprintf("%s/api/v1/partners/documents/pdf", s.urlBase)
		body        = &bytes.Buffer{}
		writer      = multipart.NewWriter(body)
	)

	part, err := writer.CreateFormFile("file", "file.pdf")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, fileReader)
	if err != nil {
		return "", err
	}

	if err := writer.WriteField("document_name", documentName); err != nil {
		return "", err
	}

	if err := writer.WriteField("teamspace_id", teamspaceID); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return "", err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var dataResponse struct {
		DocumentID string `json:"document_id"`
	}
	if err := json.NewDecoder(response.Body).Decode(&dataResponse); err != nil {
		return "", err
	}

	documentID := dataResponse.DocumentID

	return documentID, nil
}

func (s *signingServiceDoodocsSession) configureWorkflow(ctx context.Context, documentID string) error {
	var (
		url         = fmt.Sprintf("%s/api/v1/documents/%s/workflow", s.urlBase, documentID)
		jsonPayload = []byte(`{
			"workflow": {
			  "steps": [
				{
				  "index": 1,
				  "recipients": [
					{
					  "role": "anonymous_signer_rk",
					  "attrs": {
						"limit": 1
					  }
					}
				  ]
				}
			  ]
			}
		}`)
		client = &http.Client{}
	)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (s *signingServiceDoodocsSession) launchDocument(ctx context.Context, documentID string) error {
	var (
		url    = fmt.Sprintf("%s/api/v1/documents/%s/workflow/launch", s.urlBase, documentID)
		client = &http.Client{}
	)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

type recipientsResponse struct {
	Recipients []struct {
		OriginID string `json:"origin_id"`
	} `json:"recipients"`
}

func (s *signingServiceDoodocsSession) getRecipients(ctx context.Context, documentID string) (*recipientsResponse, error) {
	var (
		url    = fmt.Sprintf("%s/api/v1/documents/%s/recipients", s.urlBase, documentID)
		client = &http.Client{}
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var recipientsResponse recipientsResponse
	err = json.NewDecoder(res.Body).Decode(&recipientsResponse)
	if err != nil {
		return nil, err
	}

	return &recipientsResponse, nil
}

type recipientResponse struct {
	Recipient struct {
		Link string `json:"link"`
	} `json:"recipient"`
}

func (s *signingServiceDoodocsSession) getRecipientLink(ctx context.Context, documentID, recipientID string) (string, error) {
	var (
		url    = fmt.Sprintf("%s/api/v1/documents/%s/recipient/%s", s.urlBase, documentID, recipientID)
		client = &http.Client{}
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var recipientResponse recipientResponse
	err = json.NewDecoder(res.Body).Decode(&recipientResponse)
	if err != nil {
		return "", err
	}

	return recipientResponse.Recipient.Link, nil
}
