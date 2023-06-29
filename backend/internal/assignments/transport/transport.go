package transport

import (
	"context"
	"net/http"
	"strconv"

	"github.com/doodocs/qaztrade/backend/internal/assignments/endpoint"
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
		Limit:  int(limit),
		Offset: int(offset),
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
		Limit:  int(limit),
		Offset: int(offset),
	}, nil
}
