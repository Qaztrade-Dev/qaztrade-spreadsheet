package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/doodocs/qaztrade/backend/internal/assignments/endpoint"
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
