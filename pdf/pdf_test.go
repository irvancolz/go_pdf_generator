package pdf

import (
	"context"
	"fmt"
	"testing"
)

func TestGeneratePdf(t *testing.T) {
	toPdf := make([][]string, 0)

	headerName := []string{"no", "lorem", "ipsum", "dolor", "amet"}
	colWidth := []float64{10, 40, 40, 60, 30, 20}
	headers := GenerateTableHeaders(headerName, colWidth)

	for i := 0; i < 50; i++ {
		result := []string{fmt.Sprintf("%v", i+1), "lorem", "https://www.google.com", "kolom panjang di nomor tiga dari mungkin dengan ipsum samping", " empat"}
		toPdf = append(toPdf, result)
	}

	pdfConfig := PdfTableOptions{
		HeaderRows:  headers,
		HeaderTitle: "lorem Ipsum Dolor Sit Amet Lorem Ipsum Dolor",
	}
	ctx := context.Background()
	ExportTableToPDF(ctx, toPdf, "lorem.pdf", &pdfConfig)
	// CreateManagementFormPDF(&PdfTableOptions{PapperWidth: 600, Papperheight: 500}, CreateManagementFormTableHeader([]string{"internal", "external"}))
}

func TestGetWidth(t *testing.T) {
	header := TableHeader{
		Title: "lorem",
		Width: 123,
	}

	var child []TableHeader

	for i := 0; i < 3; i++ {
		item := TableHeader{
			Title: "hoho",
			Width: float64(i + 1),
		}

		if i == 2 {
			var gc []TableHeader
			for i := 0; i < 3; i++ {
				grandChild := TableHeader{
					Title: fmt.Sprintf("hihi %v", i+1),
					Width: 40,
				}
				gc = append(gc, grandChild)
			}

			item.Children = gc
		}

		child = append(child, item)
	}

	header.Children = child

	var tableHeaders []TableHeader
	tableHeaders = append(tableHeaders, header)
	var collumnWidth []float64

	for _, header := range tableHeaders {
		collumnWidth = append(collumnWidth, header.GetWidth(header.Children))
	}

}
