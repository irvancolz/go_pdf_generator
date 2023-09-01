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

	for i := 0; i < 100; i++ {
		result := []string{fmt.Sprintf("%v", i+1), "lorem", "https://www.google.com", "kolom panjang di nomor tiga dari mungkin dengan ipsum samping", " empat"}
		toPdf = append(toPdf, result)
	}

	pdfConfig := PdfTableOptions{
		HeaderRows: headers,
	}
	t.Log(-90 * -1)
	ctx := context.Background()
	ExportTableToPDF(ctx, toPdf, "lorem.pdf", &pdfConfig)
	// CreateManagementFormPDF(&PdfTableOptions{PapperWidth: 600, Papperheight: 500}, CreateManagementFormTableHeader([]string{"internal", "external"}))
}
