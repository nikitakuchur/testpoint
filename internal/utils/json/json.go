package json

import (
	"encoding/json"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// ReformatJson makes a given JSON pretty by adding indentation.
// Moreover, it can sort all arrays in the given JSON if you set sortArrays to true.
func ReformatJson(str string, sortArrays bool, exclude []string) string {
	var obj any
	err := json.Unmarshal([]byte(str), &obj)
	if err != nil {
		return str
	}
	return ToJson(obj, sortArrays, exclude)
}

// ToJson converts the given value into a JSON.
// Moreover, it can sort all arrays in the value before marshalling if you set sortArrays to true.
func ToJson(v any, sortArrays bool, exclude []string) string {
	newJson, ok := reformatJsonObject(v, "", sortArrays, exclude)
	if !ok {
		// everything has been excluded
		return ""
	}

	bytes, err := json.MarshalIndent(newJson, "", "  ")
	if err != nil {
		log.Fatal("failed while marshalling an object", err)
	}

	return string(bytes)
}

func reformatJsonObject(obj any, path string, sortArrays bool, exclude []string) (any, bool) {
	if isExcluded(exclude, path) {
		return nil, false
	}

	switch value := obj.(type) {
	case []interface{}:
		var arr []any
		for i, v := range value {
			v, ok := reformatJsonObject(v, path+"["+strconv.Itoa(i)+"]", sortArrays, exclude)
			if !ok {
				continue
			}
			arr = append(arr, v)
		}
		if sortArrays {
			sortArray(arr)
		}
		return arr, true
	case map[string]interface{}:
		m := make(map[string]interface{})
		for k, v := range value {
			v, ok := reformatJsonObject(v, path+"."+k, sortArrays, exclude)
			if !ok {
				continue
			}
			m[k] = v
		}
		return m, true
	default:
		return value, true
	}
}

func sortArray(array []any) {
	sort.SliceStable(array, func(i, j int) bool {
		b1, err := json.Marshal(array[i])
		if err != nil {
			log.Fatal("failed while marshalling an array", err)
		}
		b2, err := json.Marshal(array[j])
		if err != nil {
			log.Fatal("failed while marshalling an array", err)
		}
		return string(b1) < string(b2)
	})
}

func isExcluded(exclude []string, path string) bool {
	for _, e := range exclude {
		matched, _ := regexp.MatchString(wildcardToRegexp(e), path)
		if matched {
			return true
		}
	}
	return false
}

func wildcardToRegexp(pattern string) string {
	var sb strings.Builder
	for i, part := range strings.Split(pattern, "*") {
		if i > 0 {
			sb.WriteString(".*")
		}
		sb.WriteString(regexp.QuoteMeta(part))
	}
	return sb.String()
}
