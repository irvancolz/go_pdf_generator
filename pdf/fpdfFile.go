package pdf

import (
	"context"
	"fmt"
	"log"
	"math"
	"regexp"
	"time"

	"github.com/go-pdf/fpdf"
)

type TableHeader struct {
	Title    string
	Width    float64
	Children []TableHeader
}

func (t TableHeader) GetWidth(props []TableHeader) float64 {
	var result float64
	if len(t.Children) <= 0 {
		return t.Width
	}

	for _, header := range props {
		result += header.GetWidth(header.Children)
	}

	return result
}

type PdfTableOptions struct {
	// default is "Bursa Effek Indonesia"
	HeaderTitle string
	// specify each collumn name and size
	HeaderRows []TableHeader
	// "A3", "A4", "Legal", "Letter", "A5" default is "A4"
	PageSize string
	// path to logo default is globe idx
	HeaderLogo string
	// path to logo
	FooterLogo string
	// setting paper width
	PapperWidth float64
	// setting papper height
	Papperheight float64
	// setting page orientation by default "p" / "l"
	PageOrientation string
}

type fpdfPageProperties struct {
	pageLeftPadding   float64
	pageRightpadding  float64
	pageTopPadding    float64
	headerHeight      float64
	currentY          float64
	currentX          float64
	lineHeight        float64
	curColWidth       float64
	colWidthList      []float64
	totalColumn       int
	curColIndex       int
	pageHeight        float64
	currRowsheight    float64
	currPageRowHeight float64
	curRowsBgHeight   float64
	currRowsIndex     int
	tableWidth        int
	tableMarginX      float64
	currpage          int
}

func (p fpdfPageProperties) isNeedPageBreak(coord float64) bool {
	return coord > p.pageHeight-30
}

func ExportTableToPDF(c context.Context, data [][]string, filename string, props *PdfTableOptions) (string, error) {
	filenames := filename
	if filenames == "" {
		filenames = "export-to-pdf.pdf"
	}

	pdf := fpdf.NewCustom(createPageConfig(c, props))
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
	currentY := pageProps.headerHeight + 10
	lineHeight := float64(6)
	pageProps.lineHeight = lineHeight
	pageProps.currentY = currentY
	pageProps.pageHeight = pageHeight

	var columnWidth []float64
	for _, header := range props.HeaderRows {
		columnWidth = append(columnWidth, header.GetWidth(header.Children))
	}

	var totalWidth int
	for _, col := range columnWidth {
		totalWidth += int(col)
	}

	// centering the table
	tableMarginX := func() float64 {
		if totalWidth < int(pageWidth)-int(pageProps.pageLeftPadding)-int(pageProps.pageRightpadding) {
			return (pageWidth - float64(totalWidth)) / 2
		}

		return pageProps.pageLeftPadding
	}()

	pageProps.tableWidth = totalWidth
	pageProps.tableMarginX = tableMarginX
	pageProps.colWidthList = columnWidth

	currentX := tableMarginX
	pageProps.currentX = currentX
	pdf.SetLeftMargin(currentX)

	drawTableHeader(pdf, props.HeaderRows, &pageProps)
	drawTable(pdf, &pageProps, data)

	err := pdf.OutputFileAndClose(filenames)
	if err != nil {
		log.Println("failed create Pdf files :", err)
		return "", err
	}
	return filenames, nil
}

func drawTable(pdf *fpdf.Fpdf, pageProps *fpdfPageProperties, data [][]string) {
	for r, rows := range data {
		columnWidth := pageProps.colWidthList
		lineHeight := pageProps.lineHeight
		currentY := pageProps.currentY
		pageHeight := pageProps.pageHeight
		currentX := pageProps.tableMarginX
		maxColHeight := getHighestCol(pdf, columnWidth, rows)
		currRowsheight := float64(maxColHeight) * lineHeight
		currPageRowHeight := func() float64 {
			if pageProps.isNeedPageBreak(currentY+currRowsheight) && math.Floor((pageHeight-30-currentY)/lineHeight) != 0 {
				return math.Floor((pageHeight-30-currentY)/lineHeight) * lineHeight
			}
			return currRowsheight
		}()
		curRowsBgHeight := func() float64 {
			if pageProps.isNeedPageBreak(currentY+currPageRowHeight) && math.Floor((pageHeight-30-currentY)/lineHeight) != 0 {
				return math.Floor((pageHeight-30-currentY)/lineHeight) * lineHeight
			}
			return currPageRowHeight
		}()

		pageProps.currRowsheight = currRowsheight
		pageProps.currPageRowHeight = currPageRowHeight
		pageProps.curRowsBgHeight = curRowsBgHeight
		pageProps.currRowsIndex = r
		pageProps.currentX = currentX
		pageProps.currpage = pdf.PageNo()

		if pageProps.isNeedPageBreak(currentY + lineHeight) {
			pdf.AddPage()
			pdf.SetPage(pageProps.currpage + 1)
			pageProps.currpage = pdf.PageNo()
			currentY = pageProps.headerHeight + 10
			pageProps.currentY = currentY
		}
		// log.Println("current rows :", pageProps.currRowsIndex+1, " current page :", pageProps.currpage)

		pdf.SetX(currentX)
		pdf.SetY(currentY)
		pdf.SetPage(pageProps.currpage)

		drawRows(pdf, pageProps, rows)

		// do calibration on the props
		currentY = pageProps.currentY
		curRowsBgHeight = pageProps.curRowsBgHeight
		currRowsheight = pageProps.currRowsheight

		//  add height for every line added
		pageProps.currentY = func() float64 {
			if pdf.PageNo() > 1 && currentY == pageProps.headerHeight+10 && currRowsheight > curRowsBgHeight {
				return currentY + currRowsheight - curRowsBgHeight
			}
			return currentY + (float64(maxColHeight) * lineHeight)
		}()
	}
}

func drawRows(pdf *fpdf.Fpdf, pageProps *fpdfPageProperties, rows []string) {
	currentX := pageProps.tableMarginX
	currentY := pageProps.currentY
	columnWidth := pageProps.colWidthList
	lineHeight := pageProps.lineHeight
	maxColHeight := getHighestCol(pdf, columnWidth, rows)
	currRowsheight := float64(maxColHeight) * lineHeight
	curRowsBgHeight := pageProps.curRowsBgHeight

	pageProps.currRowsheight = currRowsheight
	pageProps.currentX = currentX

	pdf.SetFontStyle("")
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(240, 240, 240)

	// draw bg
	if pageProps.currRowsIndex%2 != 0 {
		pdf.SetAlpha(0, "Normal")
	}

	pdf.Rect(currentX, currentY, float64(pageProps.tableWidth), curRowsBgHeight, "F")
	pdf.SetAlpha(1, "Normal")

	pageProps.totalColumn = len(rows) - 1
	for colNumber, col := range rows {
		pageProps.curColWidth = func() float64 {
			if colNumber > len(columnWidth)-1 {
				return 20
			}

			return columnWidth[colNumber]
		}()

		pageProps.curColIndex = colNumber
		drawCell(pdf, pageProps, col)
	}
}

func drawCell(pdf *fpdf.Fpdf, pageProps *fpdfPageProperties, content string) {
	currColWidth := pageProps.curColWidth
	currentX := pageProps.currentX
	currentY := pageProps.currentY
	lineHeight := pageProps.lineHeight
	curRowsBgHeight := pageProps.curRowsBgHeight
	colNumber := pageProps.curColIndex
	totalWidth := pageProps.tableWidth
	tableMarginX := pageProps.tableMarginX

	pdf.SetY(currentY)
	pdf.SetX(currentX)

	// column border
	pdf.SetAlpha(.25, "Normal")
	pdf.Line(currentX, currentY, currentX, currentY+curRowsBgHeight)
	if colNumber == pageProps.totalColumn {
		pdf.Line(currentX+currColWidth, currentY, currentX+currColWidth, currentY+curRowsBgHeight)
	}
	pdf.SetAlpha(1, "Normal")

	splittedtext := pdf.SplitLines([]byte(content), currColWidth)
	lastRowY := currentY
	for _, text := range splittedtext {

		// reset properties when add page
		if pageProps.isNeedPageBreak(lastRowY + float64(lineHeight)) {
			pdf.AddPage()
			pdf.SetPage(pdf.PageNo() + 1)
			pageProps.currpage = pdf.PageNo()
			currentY = pageProps.headerHeight + 10
			lastRowY = currentY
			// draw bg
			if pageProps.currRowsIndex%2 != 0 {
				pdf.SetAlpha(0, "Normal")
			}
			pdf.Rect(tableMarginX, currentY, float64(totalWidth), pageProps.currRowsheight-curRowsBgHeight, "F")
			pdf.SetAlpha(1, "Normal")

			pageProps.currentY = lastRowY
		}

		pdf.SetY(lastRowY)
		pdf.SetX(currentX)

		pdf.CellFormat(currColWidth, lineHeight, string(text), "", 2, "C", false, 0, getLink(content))
		lastRowY += lineHeight
	}

	currentX += currColWidth
	pageProps.currentX = currentX
}

func drawFooter(pdf *fpdf.Fpdf) {
	pageWidth, pageHeight := pdf.GetPageSize()
	footerHeight := 10

	pdf.SetFooterFunc(func() {
		footerImgHeight := 7
		pdf.SetTextColor(0, 0, 0)
		// bottom Line
		pdf.SetFillColor(240, 240, 240)
		pdf.Rect(0, pageHeight-float64(footerHeight), pageWidth, float64(footerHeight), "F")

		// footer background
		pdf.SetDrawColor(159, 14, 15)
		pdf.SetLineWidth(.8)
		pdf.Line(0, pageHeight-float64(footerHeight), pageWidth, pageHeight-float64(footerHeight))

		// current Time
		pdf.SetFont("Times", "", 10)
		pdf.SetLeftMargin(0)
		pdf.SetY(pageHeight - float64(footerHeight))
		footerDate := func() string {
			return time.Now().Format("02/01/2006") + " - Page " + fmt.Sprintf("%v", pdf.PageNo()) + " Of " + fmt.Sprintf("%v", pdf.PageCount())
		}()
		pdf.MultiCell(50, 8, footerDate, "", "C", false)

		// app name
		pdf.SetY(pageHeight - float64(footerHeight))
		appNameWidth := 50
		pdf.SetX(pageWidth - float64(appNameWidth))
		pdf.MultiCell(float64(appNameWidth), 8, "Sistem IDX Portal", "", "C", false)

		// idx footer logo
		pdf.ImageOptions("idx-logo-2.png", (pageWidth-float64(footerImgHeight))/2, pageHeight-float64(footerHeight)+2, float64(footerImgHeight), float64(footerImgHeight), false, fpdf.ImageOptions{}, 0, "")
	})
}

func drawHeader(pdf *fpdf.Fpdf, title string, pageProps *fpdfPageProperties) {
	pageProps.headerHeight = 30
	pdf.SetHeaderFunc(func() {
		pdf.SetTextColor(0, 0, 0)
		pageWidth, pageHeight := pdf.GetPageSize()
		headerTitleWidth := pageWidth / 2
		headerTitle := title
		headerImgHeight := pageProps.headerHeight - 5
		a4PageWidth := 210
		watermarkWidth := float64(a4PageWidth) * 0.85

		// watermark
		pdf.SetAlpha(0.1, "Normal")
		pdf.ImageOptions("icon-globe-idx.png", (pageWidth-watermarkWidth)/2, (pageHeight-watermarkWidth)/2, watermarkWidth, watermarkWidth, false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
		pdf.SetAlpha(1, "Normal")

		// header bg
		pdf.SetFillColor(240, 240, 240)
		pdf.Rect(0, 0, pageWidth, float64(pageProps.headerHeight+pageProps.pageTopPadding), "F")

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

func getHighestCol(pdf *fpdf.Fpdf, colWidth []float64, data []string) int {
	result := 0

	for i, text := range data {
		currTextHeight := pdf.SplitLines([]byte(text), colWidth[i])
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
func createPageConfig(c context.Context, props *PdfTableOptions) *fpdf.InitType {
	results := fpdf.InitType{
		OrientationStr: props.getPageOrientation(c),
		UnitStr:        "mm",
		FontDirStr:     "",
		Size:           props.getPageSize(),
	}

	return &results
}

func (opt *PdfTableOptions) getPageOrientation(c context.Context) string {
	// pageOrientation := c.Query("orientation")
	// if c.Query("orientation") == "" && (opt.PageOrientation == "" || !IsContains([]string{"p", "l"}, opt.PageOrientation)) {
	// 	pageOrientation = "p"
	// 	return pageOrientation
	// }

	// if c.Query("orientation") == "" && IsContains([]string{"p", "l"}, opt.PageOrientation) {
	// 	pageOrientation = opt.PageOrientation
	// 	return pageOrientation
	// }
	return "p"
}

func (opt *PdfTableOptions) getHeaderTitle() string {
	if opt.HeaderTitle == "" {
		return "Bursa Efek Indonesia"
	}
	return opt.HeaderTitle
}

func drawTableHeader(pdf *fpdf.Fpdf, headers []TableHeader, pageProps *fpdfPageProperties) {
	if len(headers) <= 0 {
		return
	}

	currentX := pageProps.currentX
	lineHeight := pageProps.lineHeight
	currentY := pageProps.currentY

	pdf.SetFontStyle("B")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFillColor(50, 117, 168)

	// get how high this row are
	var headerTexts []string
	var headerWidths []float64
	for _, header := range headers {
		headerTexts = append(headerTexts, header.Title)
		headerWidths = append(headerWidths, header.GetWidth(header.Children))
	}

	maxColHeight := getHighestCol(pdf, headerWidths, headerTexts)
	currRowsheight := float64(maxColHeight) * lineHeight

	var totalWidth int
	for _, col := range headerWidths {
		totalWidth += int(col)
	}
	// draw bg
	pdf.Rect(currentX, currentY, float64(totalWidth), currRowsheight, "F")

	for _, header := range headers {
		pdf.SetLeftMargin(currentX)
		pdf.SetY(currentY)
		curColWidth := header.GetWidth(header.Children)

		splittedtext := pdf.SplitLines([]byte(header.Title), curColWidth)
		for _, text := range splittedtext {
			pdf.CellFormat(curColWidth, lineHeight, string(text), "", 2, "C", false, 0, getLink(header.Title))
		}

		currentX += curColWidth
		if len(header.Children) > 0 {
			currentY += lineHeight
		}

		pageProps.currentY = currentY
		drawTableHeader(pdf, header.Children, pageProps)
		pageProps.currentX = currentX
	}
	// prepare property to be used next component
	pageProps.currentY += float64(maxColHeight * int(lineHeight))
	pageProps.currentX = pageProps.tableMarginX
}

func GenerateTableHeaders(titles []string, widths []float64) []TableHeader {
	var result []TableHeader

	for i, title := range titles {
		item := TableHeader{
			Title: title,
			Width: func() float64 {
				if i >= len(widths) {
					return 20
				}
				return widths[i]
			}(),
		}
		result = append(result, item)
	}

	return result
}
