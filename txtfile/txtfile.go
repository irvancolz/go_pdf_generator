package txtfile

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func CreateTxtFromTable(data [][]string) {
	collumnWidth := getCollumnMaxWidth(data)
	file, errorCreate := os.Create("test.txt")
	if errorCreate != nil {
		log.Println("failed to create txt files :", errorCreate)
		return
	}
	defer file.Close()
	log.Println(collumnWidth)
	var beautifiedData []string

	for _, item := range data {
		var beautifiedRows strings.Builder
		for i, text := range item {
			beautifiedRows.WriteString(text)
			for space := len(text); space <= collumnWidth[i]; space++ {
				beautifiedRows.WriteString("\u0020")
			}
			beautifiedRows.WriteString("\u0020")
		}
		beautifiedData = append(beautifiedData, beautifiedRows.String())
	}

	txtFile := bufio.NewWriter(file)
	defer txtFile.Flush()

	_, errorResult := txtFile.WriteString(strings.Join(beautifiedData, "\n"))
	if errorResult != nil {
		log.Println("failed to write data to txt :", errorResult)
		return
	}

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

type TableHeaders struct {
	Name    string
	Content []TableHeaders
}

func (t TableHeaders) GetChildTotal(content []TableHeaders) int {
	var result int
	if len(content) <= 0 {
		return 1
	}

	for _, child := range content {
		result += child.GetChildTotal(child.Content)
	}
	return result
}

func CreateHeaders(headers []TableHeaders) {

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

// func getContentLength(data []interface{}) int {
// 	var result int
// 	for _, content := range data {
// 		contentType, isString := IsStringOrMap(content)
// 		if isString && contentType == "string" {
// 			result += 1
// 		} else if isString && contentType == "map" {
// 			result += getContentLength(content["children"])
// 		}
// 	}
// 	return result
// }
