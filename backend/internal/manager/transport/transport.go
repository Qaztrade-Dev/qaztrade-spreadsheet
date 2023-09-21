package transport

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
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
		applicationNoStr       = r.URL.Query().Get("application_no")
		companyNameStr         = r.URL.Query().Get("company_name")

		limit, _            = strconv.ParseUint(limitStr, 10, 0)
		offset, _           = strconv.ParseUint(offsetStr, 10, 0)
		applicationNoInt, _ = strconv.Atoi(applicationNoStr)

		filterSignedAtFrom, _  = time.Parse(time.DateOnly, filterSignedAtFromStr)
		filterSignedAtUntil, _ = time.Parse(time.DateOnly, filterSignedAtUntilStr)
	)

	return domain.GetManyInput{
		Limit:            limit,
		Offset:           offset,
		BIN:              filterBIN,
		CompensationType: filterCompensationType,
		SignedAtFrom:     filterSignedAtFrom,
		SignedAtUntil:    filterSignedAtUntil,
		CompanyName:      companyNameStr,
		ApplicationNo:    applicationNoInt,
	}, nil
}

func DecodeGetDDCard(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		applicationID = mux.Vars(r)["application_id"]
	)

	return endpoint.GetDDCardRequest{
		ApplicationID: applicationID,
	}, nil
}

func EncodeGetDDCardResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		common.EncodeError(ctx, e.Error(), w)
		return nil
	}

	var (
		resp     = response.(*endpoint.GetDDCardResponse)
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

func DecodeGetManagers(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func DecodeGrantPermissions(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		applicationID = mux.Vars(r)["application_id"]
	)

	return endpoint.GrantPermissionsRequest{
		ApplicationID: applicationID,
	}, nil
}
