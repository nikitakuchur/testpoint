package comparator

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/dop251/goja"
	"github.com/google/go-cmp/cmp"
	"github.com/sergi/go-diff/diffmatchpatch"
	"reflect"
	"testpoint/internal/sender"
)

// RespComparator is a function that performs comparison of two responses.
type RespComparator func(resp1, resp2 sender.Response) (map[string][]diffmatchpatch.Diff, error)

// NewRespComparator creates a new response comparator from the given JavaScript code.
// The script must have a function called 'compare' that accepts two responses and returns a map of diffs.
// The map of diffs can contain anything the user is interested in comparing.
// They can name keys as they want and use the 'diff' function to generate the diff.
func NewRespComparator(script string) (RespComparator, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())

	_, err := vm.RunString(script)
	if err != nil {
		return nil, errors.New("cannot run the comparator script: " + err.Error())
	}

	compare, ok := goja.AssertFunction(vm.Get("compare"))
	if !ok {
		return nil, errors.New("compare function not found")
	}

	err = vm.Set("diff", func(x, y any) []diffmatchpatch.Diff {
		if d := cmp.Diff(x, y); d != "" {
			if x != nil && y != nil && reflect.TypeOf(x).Kind() == reflect.String && reflect.TypeOf(y).Kind() == reflect.String {
				return diff(makeJsonPretty(x.(string)), makeJsonPretty(y.(string)))
			}
			return diff(toJson(x), toJson(y))
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("cannot set the diff function for js: " + err.Error())
	}

	return func(resp1, resp2 sender.Response) (map[string][]diffmatchpatch.Diff, error) {
		result, err := compare(goja.Undefined(), vm.ToValue(resp1), vm.ToValue(resp2))
		if err != nil {
			return nil, errors.New("JavaScript runtime error: " + err.Error())
		}

		return readDiffs(vm, result), nil
	}, nil
}

func toJson(v any) string {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return ""
	}
	return string(b)
}

func readDiffs(vm *goja.Runtime, v goja.Value) map[string][]diffmatchpatch.Diff {
	if v == nil || goja.IsNull(v) || goja.IsUndefined(v) {
		return nil
	}

	obj := v.ToObject(vm)

	m := make(map[string][]diffmatchpatch.Diff)
	for _, k := range obj.Keys() {
		v := obj.Get(k)
		t := v.ExportType()
		if t.Kind() == reflect.Slice && t.Elem() == reflect.TypeOf(diffmatchpatch.Diff{}) {
			diffs := obj.Get(k).Export().([]diffmatchpatch.Diff)
			if len(diffs) != 0 {
				m[k] = diffs
			}
		}
	}

	return m
}

// DefaultRespComparator is a comparator that should be used by default if a custom comparator is not provided.
func DefaultRespComparator(resp1, resp2 sender.Response) (map[string][]diffmatchpatch.Diff, error) {
	result := make(map[string][]diffmatchpatch.Diff)

	if resp1.Status != resp2.Status {
		result["status"] = diff(resp1.Status, resp2.Status)
		return result, nil
	}

	if resp1.Body != resp2.Body {
		body1, body2 := resp1.Body, resp2.Body
		if json.Valid([]byte(body1)) && json.Valid([]byte(body2)) {
			body1, body2 = makeJsonPretty(resp1.Body), makeJsonPretty(resp2.Body)
		}
		result["body"] = diff(body1, body2)
	}

	return result, nil
}

func diff(text1, text2 string) []diffmatchpatch.Diff {
	dmp := diffmatchpatch.New()
	return dmp.DiffMain(text1, text2, false)
}

func makeJsonPretty(str string) string {
	buff := bytes.Buffer{}
	err := json.Indent(&buff, []byte(str), "", "    ")
	if err != nil {
		return str
	}
	return buff.String()
}
