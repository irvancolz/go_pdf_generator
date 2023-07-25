package main

import (
	"log"
	"strings"

	"github.com/go-pdf/fpdf"
)

func loremList() []string {
	return []string{
		"Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod " +
			"tempor incididunt ut labore et dolore magna aliqua.",
		"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut " +
			"aliquip ex ea commodo consequat.",
		"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum " +
			"dolore eu fugiat nulla pariatur.",
		"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui " +
			"officia deserunt mollit anim id est laborum.",
	}
}

func main() {
	const (
		colCount = 3
		colWd    = 60.0
		marginH  = 15.0
		lineHt   = 5.5
		cellGap  = 2.0
	)
	// var colStrList [colCount]string
	type cellType struct {
		str  string
		list [][]byte
		ht   float64
	}
	var (
		cellList    [colCount]cellType
		currentCell cellType
	)

	pdf := fpdf.New("P", "mm", "A4", "") // 210 x 297
	// header := [colCount]string{"Column A", "Column B", "Column C"}
	// alignList := [colCount]string{"L", "C", "R"}
	strList := loremList()
	pdf.SetMargins(marginH, 15, marginH)
	pdf.SetFont("Arial", "", 14)
	pdf.AddPage()

	// Headers
	// pdf.SetTextColor(224, 224, 224)
	// pdf.SetFillColor(64, 64, 64)
	// for headerIndex := 0; headerIndex < colCount; headerIndex++ {
	// 	pdf.CellFormat(colWd, 10, header[headerIndex], "1", 0, "CM", true, 0, "")
	// }
	// pdf.Ln(-1)
	pdf.SetTextColor(24, 24, 24)
	pdf.SetFillColor(255, 255, 255)

	// Rows
	y := pdf.GetY()
	count := 0
	// for rows := 0; rows < 1; rows++ {
	maxHt := lineHt
	// Cell height calculation loop
	for collumn := 0; collumn < colCount; collumn++ {
		count++
		if count > len(strList) {
			count = 1
		}

		currentCell.str = strings.Join(strList[0:count], " ")
		currentCell.list = pdf.SplitLines([]byte(currentCell.str), colWd-cellGap-cellGap)
		currentCell.ht = float64(len(currentCell.list)) * lineHt

		if currentCell.ht > maxHt {
			maxHt = currentCell.ht
		}
		cellList[collumn] = currentCell
	}
	// Cell render loop
	x := marginH
	for collumnNumber := 0; collumnNumber < colCount; collumnNumber++ {
		// pdf.Rect(x, y, colWd, maxHt+cellGap+cellGap, "D")
		currentCell = cellList[collumnNumber]
		cellY := y + cellGap + (maxHt-currentCell.ht)/2

		for splittedWords := 0; splittedWords < len(currentCell.list); splittedWords++ {
			log.Println(len(currentCell.list))
			log.Println(string(currentCell.list[splittedWords]))
			pdf.SetXY(x+cellGap, cellY)
			pdf.CellFormat(colWd-cellGap-cellGap, lineHt, string(currentCell.list[splittedWords]), "", 0,
				"C", false, 0, "")
			cellY += lineHt
		}

		x += colWd
	}
	y += maxHt + cellGap + cellGap
	// }

	fileStr := "example.pdf"
	err := pdf.OutputFileAndClose(fileStr)
	if err != nil {
		log.Println(err)
		return
	}

}
