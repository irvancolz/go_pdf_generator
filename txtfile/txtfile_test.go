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
	var data []TableHeader

	for i := 0; i < 4; i++ {
		child := TableHeader{
			Name: "lorem",
		}
		data = append(data, child)
	}

	var child []TableHeader
	for i := 0; i < 12; i++ {
		childitem := TableHeader{
			Name: "lorem",
		}
		child = append(child, childitem)
	}

	withChildData := TableHeader{
		Name:     "ipsum",
		Children: child,
	}

	data = append(data, withChildData)

	var totalContent int

	for _, header := range data {
		totalContent += header.GetChildTotal(header.Children)
	}

	t.Log(totalContent)
}

func TestGetRowMaxLen(t *testing.T) {
	arrayOfArray := [][]string{
		{"apple", "banana", "orange"},
		{"cat", "dog", "elephantsloremipsum"},
		{"how", "are", "you"},
	}

	maxWidth := []int{5, 3, 4}

	modifiedArray := sliceLongText2(arrayOfArray, maxWidth)

	t.Log(modifiedArray)

}
