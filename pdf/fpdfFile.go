package pdf

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/go-pdf/fpdf"
)

type fpdfPageProperties struct {
	pageLeftPadding  float64
	pageRightpadding float64
	pageTopPadding   float64
	headerHeight     float64
	currentY         float64
}

func FpdfExport(data [][]string, props *PdfTableOptions) {

	pdf := fpdf.NewCustom(createPageConfig(props))
	pageProps := fpdfPageProperties{
		pageLeftPadding:  15,
		pageRightpadding: 15,
		pageTopPadding:   5,
	}

	pdf.SetAutoPageBreak(false, 0)
	pdf.SetFont("Arial", "", 12)
	pdf.SetMargins(pageProps.pageLeftPadding, pageProps.pageTopPadding, pageProps.pageRightpadding)

	drawHeader(pdf, props.getHeaderTitle(), &pageProps)
	drawFooter(pdf)

	pdf.SetFont("Arial", "", 12)
	pdf.AddPage()

	pageWidth, pageHeight := pdf.GetPageSize()

	columnsWidth := func() float64 {
		result := (pageWidth - pageProps.pageLeftPadding - pageProps.pageRightpadding) / float64(len(data[0]))

		if result <= 20 {
			return 20
		}

		return result
	}()

	currentY := pageProps.headerHeight + 10
	lineHeight := float64(6)

	for r, rows := range data {
		maxColHeight := getHighestCol(pdf, columnsWidth, rows)

		// reset properties when add page
		if currentY+float64(maxColHeight) > pageHeight-30 {
			pdf.AddPage()
			pdf.SetPage(pdf.PageNo() + 1)
			currentY = pageProps.headerHeight + 10
		}

		pdf.SetFontStyle("")
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFillColor(240, 240, 240)

		if r == 0 {
			pdf.SetFontStyle("B")
			pdf.SetTextColor(255, 255, 255)
			pdf.SetFillColor(50, 117, 168)
		}

		currentX := pageProps.pageLeftPadding

		if r%2 != 0 {
			pdf.SetAlpha(0, "Normal")
		}

		pdf.Rect(currentX, currentY, (pageWidth - pageProps.pageLeftPadding - pageProps.pageRightpadding), float64(maxColHeight)*lineHeight, "F")
		pdf.SetAlpha(1, "Normal")

		pdf.SetX(currentX)
		pdf.SetY(currentY)

		for _, col := range rows {
			pdf.SetY(currentY)
			pdf.SetX(currentX)

			splittedtext := pdf.SplitLines([]byte(col), columnsWidth)
			for _, text := range splittedtext {
				pdf.CellFormat(columnsWidth, lineHeight, string(text), "", 2, "C", false, 0, getLink(col))
			}
			currentX += columnsWidth
		}

		currentY += float64(maxColHeight) * lineHeight
		pageProps.currentY = currentY

	}

	err := pdf.OutputFileAndClose("lorem.pdf")
	if err != nil {
		log.Println(err)
		return
	}

}

func drawFooter(pdf *fpdf.Fpdf) {
	pageWidth, pageHeight := pdf.GetPageSize()
	footerHeight := 10

	pdf.SetFooterFunc(func() {
		footerImgHeight := 7

		// bottom Line
		pdf.SetFillColor(240, 240, 240)
		pdf.Rect(0, pageHeight-float64(footerHeight), pageWidth, float64(footerHeight), "F")

		// footer background
		pdf.SetDrawColor(159, 14, 15)
		pdf.SetLineWidth(.8)
		pdf.Line(0, pageHeight-float64(footerHeight), pageWidth, pageHeight-float64(footerHeight))

		// current Time
		pdf.SetFont("Times", "", 10)
		pdf.SetX(0)
		pdf.SetY(pageHeight - float64(footerHeight))
		footerDate := func() string {
			return time.Now().Format("02/01/2006") + " - Page " + fmt.Sprintf("%v", pdf.PageNo()) + " Of " + fmt.Sprintf("%v", pdf.PageCount())
		}()
		pdf.MultiCell(50, 8, footerDate, "", "C", false)

		// app name
		pdf.SetY(pageHeight - float64(footerHeight))
		appNameWidth := 50
		pdf.SetX(pageWidth - float64(appNameWidth))
		pdf.MultiCell(float64(appNameWidth), 8, "Sistem Portal Bursa", "", "C", false)

		// idx footer logo
		pdf.ImageOptions("idx-logo-2.png", (pageWidth-float64(footerImgHeight))/2, pageHeight-float64(footerHeight)+2, float64(footerImgHeight), float64(footerImgHeight), false, fpdf.ImageOptions{}, 0, "")
	})
}

func drawHeader(pdf *fpdf.Fpdf, title string, pageProps *fpdfPageProperties) {
	pageProps.headerHeight = 30
	pdf.SetHeaderFunc(func() {
		pageWidth, pageHeight := pdf.GetPageSize()
		headerTitleWidth := pageWidth / 2
		headerTitle := title
		headerImgHeight := pageProps.headerHeight - 5
		watermarkWidth := pageWidth * 0.85

		// header bg
		pdf.SetFillColor(240, 240, 240)
		pdf.Rect(0, 0, pageWidth, float64(pageProps.headerHeight+pageProps.pageTopPadding), "F")

		// watermark
		pdf.SetAlpha(0.1, "Normal")
		pdf.ImageOptions("icon-globe-idx.png", (pageWidth-watermarkWidth)/2, (pageHeight-watermarkWidth)/2, watermarkWidth, watermarkWidth, false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
		pdf.SetAlpha(1, "Normal")

		// header logo
		pdf.ImageOptions("icon-globe-idx.png", pageProps.pageLeftPadding, pageProps.pageTopPadding, headerImgHeight, headerImgHeight, false, fpdf.ImageOptions{}, 0, "")

		// header title
		pdf.SetLeftMargin((pageWidth - headerTitleWidth) / 2)
		pdf.SetFontSize(18)
		pdf.SetFontStyle("B")
		pdf.MultiCell(headerTitleWidth, 10, headerTitle, "", "C", false)
		pdf.SetLeftMargin(0)

		// header Bot Line
		pdf.SetDrawColor(159, 14, 15)
		pdf.SetLineWidth(.8)
		pdf.Line(0, pageProps.headerHeight+pageProps.pageTopPadding, pageWidth, pageProps.headerHeight+pageProps.pageTopPadding)

	})
}

type GinMockInterface interface {
	Query(string) string
}

// return A4 size by default
func (opt *PdfTableOptions) getPageSize() fpdf.SizeType {
	results := fpdf.SizeType{}

	results.Ht = func() float64 {
		if opt.Papperheight <= 0 {
			return 297.0
		}
		return opt.Papperheight
	}()

	results.Wd = func() float64 {
		if opt.PapperWidth <= 0 {
			return 210.0
		}
		return opt.PapperWidth
	}()

	return results
}

func (opt *PdfTableOptions) getPageOrientation(c GinMockInterface) string {
	pageOrientation := c.Query("orientation")
	if c.Query("orientation") == "" && (opt.PageOrientation == "" || !IsContains([]string{"p", "l"}, opt.PageOrientation)) {
		pageOrientation = "p"
		return pageOrientation
	}

	if c.Query("orientation") == "" && IsContains([]string{"p", "l"}, opt.PageOrientation) {
		pageOrientation = opt.PageOrientation
		return pageOrientation
	}
	return pageOrientation
}

func IsContains[T comparable](list []T, data T) bool {
	for _, item := range list {
		if item == data {
			return true
		}
	}
	return false
}

func createPageConfig(props *PdfTableOptions) *fpdf.InitType {
	results := fpdf.InitType{
		OrientationStr: props.getPageOrientation(&GinMock{}),
		UnitStr:        "mm",
		FontDirStr:     "",
		Size:           props.getPageSize(),
	}

	return &results
}

type GinMock struct {
	// data interface{}
}

func (g *GinMock) Query(props string) string {
	return "P"
}

func getHighestCol(pdf *fpdf.Fpdf, colWidth float64, data []string) int {
	result := 0

	for _, text := range data {
		currTextHeight := pdf.SplitLines([]byte(text), colWidth)
		if len(currTextHeight) > result {
			result = len(currTextHeight)
		}
	}
	return result
}

func getLink(str string) string {
	pattern := `^(https?|ftp|file):\/\/[-\w+&@#/%?=~_|!:,.;]*[-\w+&@#/%=~_|]$`
	reg := regexp.MustCompile(pattern)

	if reg.MatchString(str) {
		return str
	}

	return ""
}

type ManagementFormUserGroup struct {
	// internal I, external I etc
	Name string
	// EAB, participant. etc
	Desc string
	// checker, maker, approver
	User_Permissions []string
}

type ManagementFormTableHeader struct {
	// internal, external
	User_Type  string
	User_Group []ManagementFormUserGroup
}

func CreateManagementFormPDF(props *PdfTableOptions, headers []ManagementFormTableHeader) {
	pdf := fpdf.NewCustom(createPageConfig(props))
	pageProps := fpdfPageProperties{
		pageLeftPadding:  15,
		pageRightpadding: 15,
		pageTopPadding:   5,
	}

	pdf.SetAutoPageBreak(false, 0)
	pdf.SetFont("Arial", "", 12)
	pdf.SetMargins(pageProps.pageLeftPadding, pageProps.pageTopPadding, pageProps.pageRightpadding)

	drawHeader(pdf, props.getHeaderTitle(), &pageProps)
	drawFooter(pdf)

	pdf.SetFont("Arial", "", 12)
	pdf.AddPage()

	pageWidth, _ := pdf.GetPageSize()

	totalCollumn := func() int {
		var result int
		for _, types := range headers {
			for _, group := range types.User_Group {
				result += len(group.User_Permissions)
			}
		}

		return result
	}()

	columnsWidth := func() float64 {
		result := (pageWidth - pageProps.pageLeftPadding - pageProps.pageRightpadding) / (float64(totalCollumn) + 2)

		if result <= 20 {
			return 20
		}

		return result
	}()

	currentY := pageProps.headerHeight + 10
	currentX := pageProps.pageLeftPadding
	lineHeight := float64(6)

	// create table Header
	pdf.SetFontStyle("B")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFillColor(50, 117, 168)

	headerTableHeight := lineHeight * 4
	pdf.SetY(currentY)
	pdf.SetLeftMargin(currentX)
	pdf.CellFormat(columnsWidth, headerTableHeight, "No", "", 0, "C", true, 0, "")
	pdf.CellFormat(columnsWidth, headerTableHeight, "Nama Form", "", 0, "C", true, 0, "")
	currentX += columnsWidth + columnsWidth

	// pdf.SetX(currentX)
	for _, userType := range headers {

		currentY = pageProps.headerHeight + 10
		pdf.SetY(currentY)
		pdf.SetX(currentX)

		childCollumnTotal := func() int {
			var result int
			for _, group := range userType.User_Group {
				result += len(group.User_Permissions)
			}
			return result
		}

		pdf.CellFormat(columnsWidth*float64(childCollumnTotal()), lineHeight, userType.User_Type, "", 0, "C", true, 0, "")

		for _, group := range userType.User_Group {
			pdf.SetLeftMargin(currentX)
			currentY = pageProps.headerHeight + 10 + lineHeight
			pdf.SetY(currentY)
			pdf.CellFormat(columnsWidth*float64(len(group.User_Permissions)), lineHeight, group.Name, "", 0, "C", true, 0, "")

			currentY += lineHeight
			pdf.SetY(currentY)
			pdf.SetX(currentX)
			pdf.CellFormat(columnsWidth*float64(len(group.User_Permissions)), lineHeight, group.Desc, "", 0, "C", true, 0, "")

			for _, permission := range group.User_Permissions {
				currentY = pageProps.headerHeight + 10 + lineHeight + lineHeight + lineHeight
				pdf.SetY(currentY)
				pdf.SetX(currentX)
				pdf.CellFormat(columnsWidth, lineHeight, permission, "", 0, "C", true, 0, "")
				currentX += columnsWidth
			}
		}
	}

	err := pdf.OutputFileAndClose("manajemenform.pdf")
	if err != nil {
		log.Println("failed create Pdf files :", err)
		return
	}

}

func ExportTableToPDF(c *context.Context, data [][]string, filename string, props *PdfTableOptions) (string, error) {

	filenames := filename
	if filenames == "" {
		filenames = "export-to-pdf.pdf"
	}

	pdf := fpdf.NewCustom(createPageConfig(props))
	pageProps := fpdfPageProperties{
		pageLeftPadding:  15,
		pageRightpadding: 15,
		pageTopPadding:   5,
	}

	pdf.SetAutoPageBreak(false, 0)
	pdf.SetFont("Arial", "", 12)
	pdf.SetMargins(pageProps.pageLeftPadding, pageProps.pageTopPadding, pageProps.pageRightpadding)

	drawHeader(pdf, props.getHeaderTitle(), &pageProps)
	drawFooter(pdf)

	pdf.SetFont("Arial", "", 12)
	pdf.AddPage()

	pageWidth, pageHeight := pdf.GetPageSize()

	columnsWidth := func() float64 {
		result := (pageWidth - pageProps.pageLeftPadding - pageProps.pageRightpadding) / float64(len(data[0]))

		if result <= 20 {
			return 20
		}

		return result
	}()

	currentY := pageProps.headerHeight + 10
	lineHeight := float64(6)

	for r, rows := range data {
		maxColHeight := getHighestCol(pdf, columnsWidth, rows)

		// reset properties when add page
		if currentY+float64(maxColHeight) > pageHeight-30 {
			pdf.AddPage()
			pdf.SetPage(pdf.PageNo() + 1)
			currentY = pageProps.headerHeight + 10
		}

		pdf.SetFontStyle("")
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFillColor(240, 240, 240)

		if r == 0 {
			pdf.SetFontStyle("B")
			pdf.SetTextColor(255, 255, 255)
			pdf.SetFillColor(50, 117, 168)
		}

		currentX := pageProps.pageLeftPadding

		if r%2 != 0 {
			pdf.SetAlpha(0, "Normal")
		}

		currRowsheight := float64(maxColHeight) * lineHeight

		pdf.Rect(currentX, currentY, (pageWidth - pageProps.pageLeftPadding - pageProps.pageRightpadding), currRowsheight, "F")
		pdf.SetAlpha(1, "Normal")

		pdf.SetX(currentX)
		pdf.SetY(currentY)

		for colNumber, col := range rows {
			pdf.SetY(currentY)
			pdf.SetX(currentX)

			pdf.SetAlpha(.25, "Normal")
			pdf.Line(currentX, currentY, currentX, currentY+currRowsheight)
			if colNumber == len(rows)-1 {
				pdf.Line(pageWidth-pageProps.pageRightpadding, currentY, pageWidth-pageProps.pageRightpadding, currentY+currRowsheight)

			}
			pdf.SetAlpha(1, "Normal")

			splittedtext := pdf.SplitLines([]byte(col), columnsWidth)
			for _, text := range splittedtext {
				pdf.CellFormat(columnsWidth, lineHeight, string(text), "", 2, "C", false, 0, getLink(col))
			}
			currentX += columnsWidth
		}

		currentY += float64(maxColHeight) * lineHeight
		pageProps.currentY = currentY

	}

	err := pdf.OutputFileAndClose(filenames)
	if err != nil {
		log.Println("failed create Pdf files :", err)
		return "", err
	}
	return filenames, nil
}
