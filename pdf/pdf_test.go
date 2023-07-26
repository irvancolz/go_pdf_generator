package pdf

import (
	"fmt"
	"testing"
)

func TestGeneratePdf(t *testing.T) {
	toPdf := make([][]string, 0)

	var headers []TableHeader

	header := TableHeader{
		Title: "lorem",
		Width: 123,
	}

	header2 := TableHeader{
		Title: "lorem 2",
		Width: 40,
	}

	var child []TableHeader

	for i := 0; i < 3; i++ {
		item := TableHeader{
			Title: fmt.Sprintf("hoho index ke %v", i+1),
			Width: float64(40),
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

	// for i := 0; i < 100; i++ {
	// 	result := []string{"lorem", "https://www.google.com", "kolom panjang di nomor tiga dari samping", " empat"}
	// 	toPdf = append(toPdf, result)
	// }

	headers = append(headers, header2, header)

	pdfConfig := PdfTableOptions{
		HeaderRows:   headers,
		HeaderTitle:  "lorem Ipsum Dolor Sit Amet Lorem Ipsum Dolor",
		PapperWidth:  600,
		Papperheight: 300,
	}
	FpdfExport(toPdf, &pdfConfig)
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
