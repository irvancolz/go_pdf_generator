package csvfile

import (
	"fmt"
	"testing"
)

func TestCreateCsv(t *testing.T) {
	var toCsv [][]string
	for i := 0; i < 10; i++ {
		result := []string{"lorem", "ipsum", "dolor", fmt.Sprintf("%v", i+1)}
		toCsv = append(toCsv, result)
	}
	CreateCsv(toCsv)
}
