package transport

import (
	"context"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

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

	return endpoint.SwitchStatusRequest{
		ApplicationID: body.ApplicationID,
		StatusName:    body.StatusName,
	}, nil
}

func DecodeListSpreadsheetsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		limitStr               = r.URL.Query().Get("limit")
		offsetStr              = r.URL.Query().Get("offset")
		filterBIN              = r.URL.Query().Get("bin")
		filterCompensationType = r.URL.Query().Get("compensation_type")
		filterSignedAtFromStr  = r.URL.Query().Get("signed_at[from]")
		filterSignedAtUntilStr = r.URL.Query().Get("signed_at[until]")

		limit, _               = strconv.ParseUint(limitStr, 10, 0)
		offset, _              = strconv.ParseUint(offsetStr, 10, 0)
		filterSignedAtFrom, _  = time.Parse(time.DateOnly, filterSignedAtFromStr)
		filterSignedAtUntil, _ = time.Parse(time.DateOnly, filterSignedAtUntilStr)
	)

	return endpoint.ListSpreadsheetsRequest{
		Limit:            limit,
		Offset:           offset,
		BIN:              filterBIN,
		CompensationType: filterCompensationType,
		SignedAtFrom:     filterSignedAtFrom,
		SignedAtUntil:    filterSignedAtUntil,
	}, nil
}

func DecodeDownloadArchive(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		applicationID = mux.Vars(r)["application_id"]
	)

	return endpoint.DownloadArchiveRequest{
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

func DecodeGetDDCardResponse(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		applicationID = mux.Vars(r)["application_id"]
	)

	return endpoint.GetDDCardResponseRequest{
		ApplicationID: applicationID,
	}, nil
}

func EncodeGetDDCardResponseResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		common.EncodeError(ctx, e.Error(), w)
		return nil
	}

	var (
		resp     = response.(*endpoint.GetDDCardResponseResponse)
		httpResp = resp.HTTPResponse
	)

	defer httpResp.Body.Close()

	// copy headers from the original response
	for key, values := range httpResp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// copy the status code from the original response
	w.WriteHeader(httpResp.StatusCode)

	// copy the original response body
	_, err := io.Copy(w, httpResp.Body)
	if err != nil {
		common.EncodeError(ctx, err, w)
		return nil
	}

	return nil
}
