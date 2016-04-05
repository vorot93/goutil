package goutil

import "testing"

func TestFieldsToMap(t *testing.T) {
	var input1 = []string{"Name", "Surname", "Age"}
	var input2 = [][]string{[]string{"John", "Doe", "45"}, []string{"Jane", "Whistler", "21", "Unmarried"}}

	var fixture = []map[string]string{map[string]string{"Name": "John", "Surname": "Doe", "Age": "45"}, map[string]string{"Name": "Jane", "Surname": "Whistler", "Age": "21"}}

	var result = FieldsToMap(input1, input2)

	if !JSONcompare(fixture, result) {
		t.Error(ErrorOut(ErrMismatch, fixture, result))
		return
	}
}
