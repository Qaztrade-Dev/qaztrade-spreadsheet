package pdfservice

import (
	"bytes"
	_ "embed"
	"io"
	"text/template"

	"github.com/doodocs/qaztrade/backend/internal/sign/domain"
	"github.com/doodocs/qaztrade/backend/internal/sign/pkg/gopdf"
)

var (
	//go:embed application.tmpl
	applicationBodyBytes []byte

	//go:embed application_end.tmpl
	applicationEndBodyBytes []byte

	headerText = "Заявка на получение возмещения части затрат субъектов индустриально-инновационной деятельности"
)

type PDFService struct {
	applicationTemplate    *template.Template
	applicationEndTemplate *template.Template
}

var _ domain.PDFService = (*PDFService)(nil)

func NewPDFService() (*PDFService, error) {
	applicationTemplate, err := template.New("").Parse(string(applicationBodyBytes))
	if err != nil {
		return nil, err
	}

	applicationEndTemplate, err := template.New("").Parse(string(applicationEndBodyBytes))
	if err != nil {
		return nil, err
	}

	svc := &PDFService{
		applicationTemplate:    applicationTemplate,
		applicationEndTemplate: applicationEndTemplate,
	}
	return svc, nil
}

func (s *PDFService) Create(application *domain.Application, attachments []io.ReadSeeker) (*bytes.Buffer, error) {
	pdfFirstPart, err := s.createFirstPart(application, attachments)
	if err != nil {
		return nil, err
	}

	pdfFull, err := s.createSecondPart(application, []io.ReadSeeker{bytes.NewReader(pdfFirstPart.Bytes())})
	if err != nil {
		return nil, err
	}

	return pdfFull, nil
}

func (s *PDFService) createFirstPart(application *domain.Application, attachments []io.ReadSeeker) (*bytes.Buffer, error) {
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
	pdf.WithPDFsAfter(attachments...)

	pdfBuffer, err := pdf.Output()
	if err != nil {
		return nil, err
	}

	return pdfBuffer, nil
}

func (s *PDFService) createSecondPart(application *domain.Application, attachments []io.ReadSeeker) (*bytes.Buffer, error) {
	var temp bytes.Buffer

	if err := s.applicationEndTemplate.Execute(&temp, application); err != nil {
		return nil, err
	}

	pdf, err := gopdf.NewPDF()
	if err != nil {
		return nil, err
	}

	page1 := gopdf.NewPage(
		&gopdf.Text{
			Font:     gopdf.FontRegular,
			FontSize: 10,
			Align:    "LT",
			Text:     temp.String(),
		},
	)

	pdf.WithPages(page1)
	pdf.WithPDFsBefore(attachments...)

	pdfBuffer, err := pdf.Output()
	if err != nil {
		return nil, err
	}

	return pdfBuffer, nil
}
