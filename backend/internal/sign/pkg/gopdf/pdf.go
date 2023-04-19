package gopdf

import (
	"bytes"
	_ "embed"
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/phpdave11/gofpdf"
)

const (
	PageType         = "A4"
	PageOrientation  = "P"
	PageUnits        = "mm"
	PageWidth        = 210
	PageHeight       = 297
	PageTopMargin    = 10
	PageBottomMargin = 10
	PageLeftMargin   = 30
	PageRightMargin  = 10
	ContentMaxWidth  = PageWidth - PageLeftMargin - PageRightMargin
	ContentMaxHeight = PageHeight - PageTopMargin - PageBottomMargin
	HeaderHeight     = 10
	FooterHeight     = 10

	FontRegular            = "LiberationSans-Regular"
	FontBold               = "LiberationSans-Bold"
	FontItalic             = "LiberationSans-Italic"
	FontBoldItalic         = "LiberationSans-BoldItalic"
	FontSuisseIntl         = "SuisseIntl-Regular"
	FontSuisseIntlBold     = "SuisseIntl-Bold"
	FontSuisseIntlSemiBold = "SuisseIntl-SemiBold"
)

var (
	//go:embed fonts/LiberationSans-Regular.ttf
	embeddedFontRegular []byte

	//go:embed fonts/LiberationSans-Bold.ttf
	embeddedFontBold []byte

	//go:embed fonts/LiberationSans-Italic.ttf
	embeddedFontItalic []byte

	//go:embed fonts/LiberationSans-BoldItalic.ttf
	embeddedFontBoldItalic []byte

	//go:embed fonts/SuisseIntl/SuisseIntl-Regular.ttf
	fontSuisseIntl []byte

	//go:embed fonts/SuisseIntl/SuisseIntl-Bold.ttf
	fontSuisseIntlBold []byte

	//go:embed fonts/SuisseIntl/SuisseIntl-SemiBold.ttf
	fontSuisseIntlSemiBold []byte
)

type PDF struct {
	fpdf       *gofpdf.Fpdf
	pages      []*Page
	isRendered bool
	beforePDFs []io.ReadSeeker
	afterPDFs  []io.ReadSeeker
}

func NewPDF() (*PDF, error) {
	fpdf, err := newFpdf()
	if err != nil {
		return nil, err
	}

	pdf := &PDF{
		fpdf:  fpdf,
		pages: make([]*Page, 0),
	}

	return pdf, nil
}

func newFpdf() (pdf *gofpdf.Fpdf, err error) {
	pdf = gofpdf.New(PageOrientation, PageUnits, PageType, "")

	// Fpdf by default sets PDF version to "1.3" and not always bumps it when uses newer features.
	// Adding an empty layer bumps the version to "1.5" thus increasing compliance with the standard.
	pdf.AddLayer("Empty", false)

	pdf.AddUTF8FontFromBytes(FontRegular, "", embeddedFontRegular)
	pdf.AddUTF8FontFromBytes(FontBold, "", embeddedFontBold)
	pdf.AddUTF8FontFromBytes(FontItalic, "", embeddedFontItalic)
	pdf.AddUTF8FontFromBytes(FontBoldItalic, "", embeddedFontBoldItalic)
	pdf.AddUTF8FontFromBytes(FontSuisseIntl, "", fontSuisseIntl)
	pdf.AddUTF8FontFromBytes(FontSuisseIntlBold, "", fontSuisseIntlBold)
	pdf.AddUTF8FontFromBytes(FontSuisseIntlSemiBold, "", fontSuisseIntlSemiBold)

	// Fpdf margins are used only on Info Block pages, configure them with header and footer height to utilize auto page break
	pdf.SetMargins(PageLeftMargin, PageTopMargin+HeaderHeight, PageRightMargin)
	pdf.SetAutoPageBreak(true, PageBottomMargin+FooterHeight)

	if err := pdf.Error(); err != nil {
		return nil, err
	}

	return pdf, nil
}

type Page struct {
	Texts []*Text
}

func NewPage(texts ...*Text) *Page {
	return &Page{
		Texts: texts,
	}
}

type Text struct {
	Font     string
	FontSize float64
	Align    string // "LT", "CT"
	Text     string
}

func (p *Page) WithTexts(texts ...*Text) {
	p.Texts = append(p.Texts, texts...)
}

func (t *Text) Render(fpdf *gofpdf.Fpdf) {
	fpdf.SetFont(t.Font, "", t.FontSize)
	fpdf.MultiCell(ContentMaxWidth, 4, t.Text, "", t.Align, false)
}

func (p *Page) Render(fpdf *gofpdf.Fpdf) {
	for _, text := range p.Texts {
		text.Render(fpdf)
	}
}

func (p *PDF) WithPDFsBefore(pdfs ...io.ReadSeeker) {
	p.beforePDFs = append(p.beforePDFs, pdfs...)
}

func (p *PDF) WithPDFsAfter(pdfs ...io.ReadSeeker) {
	p.afterPDFs = append(p.afterPDFs, pdfs...)
}

func (p *PDF) WithPages(pages ...*Page) {
	p.pages = append(p.pages, pages...)
}

func (p *PDF) Render() error {
	for _, page := range p.pages {
		p.fpdf.AddPage()
		page.Render(p.fpdf)

		if err := p.fpdf.Error(); err != nil {
			return err
		}
	}
	p.isRendered = true
	return nil
}

func (p *PDF) Output() (*bytes.Buffer, error) {
	if !p.isRendered {
		if err := p.Render(); err != nil {
			return nil, err
		}
	}

	var buffer = &bytes.Buffer{}

	if err := p.fpdf.Output(buffer); err != nil {
		return nil, err
	}

	if len(p.beforePDFs) > 0 {
		return p.mergeBefore(buffer)
	}

	if len(p.afterPDFs) > 0 {
		return p.mergeAfter(buffer)
	}

	return buffer, nil
}

func (p *PDF) mergeBefore(inputBuffer *bytes.Buffer) (*bytes.Buffer, error) {
	var (
		buffer     = &bytes.Buffer{}
		readSeeker = bytes.NewReader(inputBuffer.Bytes())
		rss        = append(p.beforePDFs, readSeeker)
	)

	if err := api.Merge(rss, buffer, nil); err != nil {
		return nil, err
	}

	return buffer, nil
}

func (p *PDF) mergeAfter(inputBuffer *bytes.Buffer) (*bytes.Buffer, error) {
	var (
		buffer     = &bytes.Buffer{}
		readSeeker = bytes.NewReader(inputBuffer.Bytes())
		rss        = append([]io.ReadSeeker{readSeeker}, p.afterPDFs...)
	)

	if err := api.Merge(rss, buffer, nil); err != nil {
		return nil, err
	}

	return buffer, nil
}
