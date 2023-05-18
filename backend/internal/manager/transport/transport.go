package transport

import (
	"context"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/doodocs/qaztrade/backend/internal/manager/endpoint"
	"github.com/gorilla/mux"
)

func DecodeSwitchStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		ApplicationID string `json:"application_id"`
		StatusName    string `json:"status_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	tokenString := extractHeaderToken(r)

	return endpoint.SwitchStatusRequest{
		UserToken:     tokenString,
		ApplicationID: body.ApplicationID,
		StatusName:    body.StatusName,
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

func DecodeListSpreadsheetsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		limitStr    = r.URL.Query().Get("limit")
		offsetStr   = r.URL.Query().Get("offset")
		tokenString = extractHeaderToken(r)

		limit, _  = strconv.ParseUint(limitStr, 10, 0)
		offset, _ = strconv.ParseUint(offsetStr, 10, 0)
	)

	return endpoint.ListSpreadsheetsRequest{
		UserToken: tokenString,
		Limit:     limit,
		Offset:    offset,
	}, nil
}

func DecodeDownloadArchive(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		applicationID = mux.Vars(r)["application_id"]
		tokenString   = extractHeaderToken(r)
	)

	return endpoint.DownloadArchiveRequest{
		UserToken:     tokenString,
		ApplicationID: applicationID,
	}, nil
}

func EncodeDownloadArchiveResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		common.EncodeError(ctx, e.Error(), w)
		return nil
	}

	resp := response.(*endpoint.DownloadArchiveResponse)
	defer resp.RemoveFunc()
	defer resp.ArchiveReader.Close()

	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext("архив.zip")))
	w.Header().Set("Content-Disposition", "attachment; filename=\"архив.zip\"")

	_, err := io.Copy(w, resp.ArchiveReader)
	if err != nil {
		common.EncodeError(ctx, err, w)
		return nil
	}

	return nil
}
