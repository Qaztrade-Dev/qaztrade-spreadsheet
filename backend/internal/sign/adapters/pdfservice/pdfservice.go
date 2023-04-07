package main

import (
	"bytes"
	_ "embed"
	"io"
	"text/template"

	"github.com/doodocs/qaztrade/backend/internal/sign/pkg/gopdf"
	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
)

var (
	//go:embed application.tmpl
	applicationBodyBytes []byte

	headerText = "Заявка на получение возмещения части затрат субъектов индустриально-инновационной деятельности"
)

type PDFService struct {
	applicationTemplate *template.Template
}

func NewPDFService() (*PDFService, error) {
	applicationTemplate, err := template.New("").Parse(string(applicationBodyBytes))
	if err != nil {
		return nil, err
	}

	svc := &PDFService{
		applicationTemplate: applicationTemplate,
	}
	return svc, nil
}

func (s *PDFService) Create(application *domain.Application, attachments []io.ReadSeeker) (*bytes.Buffer, error) {
	var temp bytes.Buffer

	if err := s.applicationTemplate.Execute(&temp, application); err != nil {
		return nil, err
	}

	pdf, err := gopdf.NewPDF()
	if err != nil {
		return nil, err
	}

	page1 := gopdf.NewPage(
		&gopdf.Text{
			Font:     gopdf.FontBold,
			FontSize: 10,
			Align:    "CT",
			Text:     headerText,
		},
		&gopdf.Text{
			Font:     gopdf.FontRegular,
			FontSize: 10,
			Align:    "LT",
			Text:     temp.String(),
		},
	)

	pdf.WithPages(page1)
	pdf.WithPDFs(attachments...)

	pdfBuffer, err := pdf.Output()
	if err != nil {
		return nil, err
	}

	return pdfBuffer, nil
}
