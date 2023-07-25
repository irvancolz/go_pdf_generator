package txtfile

import (
	"fmt"
	"testing"
)

func TestCreateTxt(t *testing.T) {
	var toTxt [][]string
	for i := 0; i < 100; i++ {
		result := []string{"lorem", fmt.Sprintf("%v", i+1), "ipsum", "dolor"}
		if i == 12 {
			result[2] = "ipsumametasdkfjnakfnanfn"
		}
		toTxt = append(toTxt, result)
	}
	CreateTxtFromTable(toTxt)
}

func TestIsStrings(t *testing.T) {
	_, isString := IsStringOrMap("lorem")
	if !isString {
		t.Error("the function should return true")
	}
}

func TestIsNotStrings(t *testing.T) {
	_, isString := IsStringOrMap(1)
	if isString {
		t.Error("the function should return false")
	}
}

func TestIsStrings2(t *testing.T) {
	types, isString := IsStringOrMap("lorem")
	if !isString || types != "string" {
		t.Error("the function should return true and type of 'string'")
	}

}
func TestIsMap(t *testing.T) {
	data := map[string]interface{}{"title": 1}

	types, isMap := IsStringOrMap(data)
	if !isMap || types != "map" {
		t.Error("the function should return true and type of 'map'")
	}
}

func TestGetContentlength(t *testing.T) {
	var data []TableHeaders

	for i := 0; i < 4; i++ {
		child := TableHeaders{
			Name: "lorem",
		}
		data = append(data, child)
	}

	var child []TableHeaders
	for i := 0; i < 12; i++ {
		childitem := TableHeaders{
			Name: "lorem",
		}
		child = append(child, childitem)
	}

	withChildData := TableHeaders{
		Name:    "ipsum",
		Content: child,
	}

	data = append(data, withChildData)

	var totalContent int

	for _, header := range data {
		totalContent += header.GetChildTotal(header.Content)
	}

	t.Log(totalContent)
}
