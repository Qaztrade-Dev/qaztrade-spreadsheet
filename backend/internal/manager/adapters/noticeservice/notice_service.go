package noticeservice

import (
	"bytes"
	_ "embed"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
	"github.com/lukasjarosch/go-docx"
)

//go:embed docs_temp.docx
var docxFile []byte

type NoticeService struct {
	docxTemplate []byte
}

func NewNoticeService() (*NoticeService, error) {
	svc := &NoticeService{
		docxTemplate: docxFile,
	}
	return svc, nil
}

var _ domain.NoticeService = (*NoticeService)(nil)

func (s *NoticeService) Create(revision *domain.Revision) (*bytes.Buffer, error) {
	replaceMap := docx.PlaceholderMap{
		"CompanyName":     revision.To,
		"CompanyAddress":  revision.Address,
		"ApplicationNum":  revision.No,
		"ApplicationDate": revision.CreatedAt.Format("2006-01-02"),
		"Remarks":         revision.Remarks,
	}
	doc, err := docx.OpenBytes(docxFile)
	if err != nil {
		return nil, err
	}

	err = doc.ReplaceAll(replaceMap)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = doc.Write(&buf)
	if err != nil {
		return nil, err
	}
	return &buf, nil

}
