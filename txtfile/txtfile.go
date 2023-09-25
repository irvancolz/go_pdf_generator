package txtfile

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func CreateTxtFromTable(data [][]string, columnWidth []int) {
	// collumnWidth := getCollumnMaxWidth(data)
	file, errorCreate := os.Create("test.txt")
	if errorCreate != nil {
		log.Println("failed to create txt files :", errorCreate)
		return
	}
	defer file.Close()

	// beautifiedData := drawTxtTable(data, collumnWidth)
	txtFile := bufio.NewWriter(file)
	defer txtFile.Flush()
	beautifiedData := formatRowsData(data, columnWidth)
	withStylingData := drawTxtTable(beautifiedData, columnWidth)

	_, errorResult := txtFile.WriteString(strings.Join(withStylingData, "\n"))

	if errorResult != nil {
		log.Println("failed to write data to txt :", errorResult)
		return
	}

}

func getRowMaxContent(data []string) int {
	var result int

	for _, word := range data {
		if len(word) >= result {
			result = len(word)
		}
	}

	return result
}

func getCollumnMaxWidth(data [][]string) []int {
	var result []int
	var tableCollumnContent [][]string

	for i := 0; i < len(data[0]); i++ {
		tableCollumnContent = append(tableCollumnContent, []string{})
	}

	for _, rows := range data {
		for d, text := range rows {
			tableCollumnContent[d] = append(tableCollumnContent[d], text)
		}
	}

	for _, rows := range tableCollumnContent {
		var max int
		for _, text := range rows {
			if len(text) > max {
				max = len(text)
			}
		}
		result = append(result, max)
	}
	return result
}

type TableHeader struct {
	Name     string
	Width    float64
	Children []TableHeader
}

func (t TableHeader) GetChildTotal(content []TableHeader) int {
	var result int
	if len(content) <= 0 {
		return 1
	}

	for _, child := range content {
		result += child.GetChildTotal(child.Children)
	}
	return result
}

func CreateHeaders(headers []TableHeader) {

}

func IsStringOrMap(target interface{}) (string, bool) {
	switch target.(type) {
	case string:
		{
			return "string", true
		}
	case map[string]interface{}:
		{
			return "map", true
		}
	default:
		{
			return "", false
		}
	}
}

func drawTxtTable(data [][]string, columnWidth []int) []string {
	var result []string
	for j, item := range data {
		var beautifiedRows strings.Builder
		for i, text := range item {
			beautifiedRows.WriteString("| " + removeEscEnter(text))
			for space := len(text); space <= columnWidth[i]; space++ {
				beautifiedRows.WriteString("\u0020")
			}
			if i == len(item)-1 {
				beautifiedRows.WriteString("|")
			}
		}
		// the 2 pipe "|" before and after the text and the space " " after the first pipe
		totalStylingCharacter := 3
		// header border bottom
		if j == 0 {
			beautifiedRows.WriteString("\n")
			for i := 0; i < len(item); i++ {
				for char := 0; char <= columnWidth[i]+totalStylingCharacter; char++ {
					beautifiedRows.WriteString("-")
				}
			}
		}

		result = append(result, beautifiedRows.String())
	}
	return result
}

func checkOverlappedText(data [][]string, widths []int) bool {
	var result bool

	for _, line := range data {

		for i, word := range line {
			maxWordLen := func() int {
				if i >= len(widths) {
					return 20
				}

				return widths[i]
			}()
			if len(word) > maxWordLen {
				return true
			}
		}
	}

	return result
}

func checkAndModifyArray(arr [][]string) [][]string {
	var resultArr [][]string

	for _, subArr := range arr {
		var modifiedSubArr []string

		for _, str := range subArr {
			if len(str) > 5 {
				modifiedSubArr = append(modifiedSubArr, str[:5])
				remainingStr := str[5:]
				for len(remainingStr) > 5 {
					modifiedSubArr = append(modifiedSubArr, remainingStr[:5])
					remainingStr = remainingStr[5:]
				}
				if len(remainingStr) > 0 {
					modifiedSubArr = append(modifiedSubArr, remainingStr)
				}
			} else {
				modifiedSubArr = append(modifiedSubArr, str)
			}
		}

		resultArr = append(resultArr, modifiedSubArr)
	}

	return resultArr
}

func sliceLongText2(data [][]string, maxWidths []int) [][]string {
	var result [][]string
	for line := 0; line < len(data); line++ {
		copyOfLine := data[line]
		newLineMock := make([]string, len(copyOfLine))
		for wordIdx, word := range data[line] {
			maxWordLen := func() int {
				if wordIdx >= len(maxWidths) {
					return 20
				}

				return maxWidths[wordIdx]
			}()
			if len(word) > maxWordLen {
				slicedWord := string([]byte(word)[:maxWordLen])
				remainingWord := string([]byte(word)[maxWordLen:])
				copyOfLine[wordIdx] = slicedWord
				newLineMock[wordIdx] = remainingWord
			}
		}
		result = append(result, copyOfLine)
		if getRowMaxContent(newLineMock) != 0 {
			result = append(result, newLineMock)
		}
	}

	if !checkOverlappedText(result, maxWidths) {
		return result
	}

	return sliceLongText2(result, maxWidths)
}

func removeEscEnter(s string) string {
	return strings.Join(strings.Split(strings.ReplaceAll(strconv.Quote(s), `"`, ""), `\n`), " ")
}

func formatRowsData(data [][]string, maxWidths []int) [][]string {

	var result [][]string
	for line := 0; line < len(data); line++ {
		copyOfLine := data[line]
		newLineMock := make([]string, len(copyOfLine))
		for wordIdx, word := range data[line] {
			maxWordLen := func() int {
				if wordIdx >= len(maxWidths) {
					return 20
				}

				return maxWidths[wordIdx]
			}()
			if len(word) > maxWordLen {
				slicedWord := string([]byte(word)[:maxWordLen])
				remainingWord := string([]byte(word)[maxWordLen:])
				copyOfLine[wordIdx] = slicedWord
				newLineMock[wordIdx] = remainingWord
			}
		}
		result = append(result, copyOfLine)
		if getRowMaxContent(newLineMock) != 0 {
			result = append(result, newLineMock)
		}
	}

	if !checkOverlappedText(result, maxWidths) {
		return result
	}

	return formatRowsData(result, maxWidths)

}
