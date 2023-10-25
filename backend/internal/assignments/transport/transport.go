package transport

import (
	"context"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"github.com/doodocs/qaztrade/backend/internal/assignments/endpoint"
	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/gorilla/mux"
)

func DecodeGetAssignmentsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		limitStr          = r.URL.Query().Get("limit")
		offsetStr         = r.URL.Query().Get("offset")
		assigneeIDStr     = r.URL.Query().Get("assignee_id")
		companyNameStr    = r.URL.Query().Get("company_name")
		applicationNoStr  = r.URL.Query().Get("application_no")
		assignmentTypeStr = r.URL.Query().Get("assignment_type")

		limit, _            = strconv.ParseUint(limitStr, 10, 0)
		offset, _           = strconv.ParseUint(offsetStr, 10, 0)
		applicationNoInt, _ = strconv.Atoi(applicationNoStr)

		applicationNo  *int
		assigneeID     *string
		companyName    *string
		assignmentType *string
	)

	if assigneeIDStr != "" {
		assigneeID = &assigneeIDStr
	}

	if companyNameStr != "" {
		companyName = &companyNameStr
	}

	if applicationNoStr != "" {
		applicationNo = &applicationNoInt
	}

	if assignmentTypeStr != "" {
		assignmentType = &assignmentTypeStr
	}

	return domain.GetManyInput{
		Limit:          limit,
		Offset:         offset,
		AssigneeID:     assigneeID,
		CompanyName:    companyName,
		ApplicationNo:  applicationNo,
		AssignmentType: assignmentType,
	}, nil
}

func DecodeGetUserAssignmentsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		limitStr          = r.URL.Query().Get("limit")
		offsetStr         = r.URL.Query().Get("offset")
		companyNameStr    = r.URL.Query().Get("company_name")
		applicationNoStr  = r.URL.Query().Get("application_no")
		assignmentTypeStr = r.URL.Query().Get("assignment_type")

		limit, _            = strconv.ParseUint(limitStr, 10, 0)
		offset, _           = strconv.ParseUint(offsetStr, 10, 0)
		applicationNoInt, _ = strconv.Atoi(applicationNoStr)

		applicationNo  *int
		companyName    *string
		assignmentType *string
	)

	if companyNameStr != "" {
		companyName = &companyNameStr
	}

	if applicationNoStr != "" {
		applicationNo = &applicationNoInt
	}

	if assignmentTypeStr != "" {
		assignmentType = &assignmentTypeStr
	}

	return domain.GetManyInput{
		Limit:          limit,
		Offset:         offset,
		CompanyName:    companyName,
		ApplicationNo:  applicationNo,
		AssignmentType: assignmentType,
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

func DecodeGetArchiveRequest(_ context.Context, r *http.Request) (interface{}, error) {
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
		return err
	}

	return nil
}

func DecodeCheckAssignmentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		assignmentIDStr = mux.Vars(r)["assignment_id"]
		assignmentID, _ = strconv.ParseUint(assignmentIDStr, 10, 0)
	)

	return endpoint.CheckAssignmentRequest{
		AssignmentID: assignmentID,
	}, nil
}

func DecodeEnqueueAssignmentsRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func DecodeRedistributeAssignmentsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		assignmentType = r.URL.Query().Get("assignment_type")
	)
	return endpoint.RedistributeAssignmentsRequest{
		AssignmentType: assignmentType,
	}, nil
}

func DecodeSendNotice(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		assignmentIDStr = mux.Vars(r)["assignment_id"]
		assignmentID, _ = strconv.ParseUint(assignmentIDStr, 10, 0)
	)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return nil, err
	}

	fileReader, header, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	var (
		fileSize = header.Size
	)

	return endpoint.SendNoticeRequest{
		AssignmentID: assignmentID,
		FileReader:   fileReader,
		FileSize:     fileSize,
		FileName:     header.Filename,
	}, nil
}

func DecodeRespondNotice(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		applicationID  = mux.Vars(r)["application_id"]
		assignmentType = mux.Vars(r)["assignment_type"]
	)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return nil, err
	}

	fileReader, header, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	var (
		fileSize = header.Size
	)

	return endpoint.RespondNoticeRequest{
		ApplicationID:  applicationID,
		AssignmentType: assignmentType,
		FileReader:     fileReader,
		FileSize:       fileSize,
		FileName:       header.Filename,
	}, nil
}
