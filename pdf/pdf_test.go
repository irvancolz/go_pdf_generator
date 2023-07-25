package pdf

import (
	"testing"
)

func TestGeneratePdf(t *testing.T) {
	toPdf := make([][]string, 0)

	for i := 0; i < 100; i++ {
		result := []string{"lorem", "https://www.google.com", "kolom panjang di nomor tiga dari samping", " empat"}
		toPdf = append(toPdf, result)
	}

	FpdfExport(toPdf, &PdfTableOptions{HeaderTitle: "lorem Ipsum Dolor Sit Amet Lorem Ipsum Dolor"})
	// CreateManagementFormPDF(&PdfTableOptions{PapperWidth: 600, Papperheight: 500}, CreateManagementFormTableHeader([]string{"internal", "external"}))
}

// func TestGetLink(t *testing.T) {
// 	stringsLink := []string{
// 		"example.com",
// 		"www.example.com",
// 		"http://example.com",
// 		"https://www.example.com",
// 		"ftp://example.com",
// 		"invalid-url",
// 	}
// 	for _, link := range stringsLink {
// 		t.Log("link untuk ", link, "adalah :", getLink(link))
// 	}
// }
