package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/doodocs/qaztrade/backend/internal/assignments/endpoint"
	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/gorilla/mux"
)

func DecodeGetAssignmentsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		limitStr  = r.URL.Query().Get("limit")
		offsetStr = r.URL.Query().Get("offset")
		userIDStr = r.URL.Query().Get("user_id")

		limit, _  = strconv.ParseUint(limitStr, 10, 0)
		offset, _ = strconv.ParseUint(offsetStr, 10, 0)
		userID    *string
	)

	if userIDStr != "" {
		userID = &userIDStr
	}

	return endpoint.GetAssignmentsRequest{
		Limit:  limit,
		Offset: offset,
		UserID: userID,
	}, nil
}

func DecodeGetUserAssignmentsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		limitStr  = r.URL.Query().Get("limit")
		offsetStr = r.URL.Query().Get("offset")

		limit, _  = strconv.ParseUint(limitStr, 10, 0)
		offset, _ = strconv.ParseUint(offsetStr, 10, 0)
	)

	return endpoint.GetUserAssignmentsRequest{
		Limit:  limit,
		Offset: offset,
	}, nil
}

func DecodeCreateBatchRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func DecodeChangeAssigneeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		assignmentIDStr = mux.Vars(r)["assignment_id"]
		assignmentID, _ = strconv.ParseUint(assignmentIDStr, 10, 0)
	)

	var body struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return endpoint.ChangeAssigneeRequest{
		UserID:       body.UserID,
		AssignmentID: assignmentID,
	}, nil
}

func DecodeGetArchive(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		assignmentIDStr = mux.Vars(r)["assignment_id"]
		assignmentID, _ = strconv.ParseUint(assignmentIDStr, 10, 0)
	)

	return endpoint.GetArchiveRequest{
		AssignmentID: assignmentID,
	}, nil
}

func EncodeGetArchiveResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(common.Errorer); ok && e.Error() != nil {
		common.EncodeError(ctx, e.Error(), w)
		return nil
	}

	resp := response.(*endpoint.GetArchiveResponse)
	defer resp.RemoveFunc()
	defer resp.ArchiveReader.Close()

	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext("архив.zip")))
	w.Header().Set("Content-Disposition", "attachment; filename=\"архив.zip\"")

	_, err := io.Copy(w, resp.ArchiveReader)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil
}
