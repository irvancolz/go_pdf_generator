package excel

import (
	"fmt"
	"testing"
)

type Lorem struct {
	Ipsum  string
	Dolor  string
	Sit    string
	Amet   int
	Kokoro string
}

func TestConvertToArray(t *testing.T) {
	var list []interface{}

	config := ExportToExcelConfig{
		HeaderText: []string{"CONTACT PERSON ANGGOTA BURSA / PARTISIPAN / PJ SPPA / DU", "Kode :	BUDS", "Nama Perusahaan : 	BUMI DAYA SEKURITAS"},
	}

	for i := 0; i < 10; i++ {
		data := Lorem{
			Ipsum:  fmt.Sprintf("ipsum koewanf snadfkn %d", i),
			Dolor:  fmt.Sprintf("dolor %d", i),
			Sit:    fmt.Sprintf("sit %d", i),
			Amet:   i,
			Kokoro: fmt.Sprintf("kokoro %d", i),
		}
		list = append(list, data)
	}

	var toExcelData [][]string
	for _, item := range list {
		// result := MapToArray(item, []string{"ipsum", "dolor", "amet"})
		result := StructToArray(item, []string{"Dolor", "Amet", "Ipsum", "Kokoro"})
		toExcelData = append(toExcelData, result)
	}
	result, errExport := config.ExportTableToExcel("lorem", toExcelData)
	if errExport != nil {
		return
	}
	t.Log(result)
}

func TestReadFile(t *testing.T) {
	excelValue := ReadFileExcel("../lorem.xlsx")
	if excelValue != nil {
		t.Log(excelValue)
	}
}
