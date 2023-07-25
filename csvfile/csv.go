package csvfile

import (
	"encoding/csv"
	"log"
	"os"
)

func CreateCsv(data [][]string) {
	file, errorFile := os.Create("testing.csv")
	if errorFile != nil {
		log.Println("failed to create csv file", errorFile)
		return
	}
	defer file.Close()

	csvFile := csv.NewWriter(file)
	defer csvFile.Flush()

	csvFile.WriteAll(data)
}
