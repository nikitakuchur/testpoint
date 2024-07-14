package json

import (
	"encoding/json"
	"sort"
)

// ReformatJson makes a given JSON pretty by adding indentation.
// Moreover, it can sort all arrays in the given JSON if you set sortArrays to true.
func ReformatJson(str string, sortArrays bool) string {
	var obj any
	err := json.Unmarshal([]byte(str), &obj)
	if err != nil {
		panic(err)
	}
	return ToJson(obj, sortArrays)
}

// ToJson converts the given value into a JSON.
// Moreover, it can sort all arrays in the value before marshalling if you set sortArrays to true.
func ToJson(v any, sortArrays bool) string {
	if sortArrays {
		sortJsonArrays(v)
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func sortJsonArrays(obj any) {
	switch value := obj.(type) {
	case []interface{}:
		for _, v := range value {
			sortJsonArrays(v)
		}
		sort.SliceStable(value, func(i, j int) bool {
			b1, err := json.Marshal(value[i])
			if err != nil {
				panic(err)
			}
			b2, err := json.Marshal(value[j])
			if err != nil {
				panic(err)
			}
			return string(b1) < string(b2)
		})
	case map[string]interface{}:
		for _, v := range value {
			sortJsonArrays(v)
		}
	}
}
