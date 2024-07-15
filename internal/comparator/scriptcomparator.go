package comparator

import (
	"errors"
	"github.com/dop251/goja"
	"github.com/sergi/go-diff/diffmatchpatch"
	"reflect"
	"testpoint/internal/diff"
	"testpoint/internal/sender"
	jsonutils "testpoint/internal/utils/json"
)

type ScriptComparator struct {
	runtime *goja.Runtime
	compare goja.Callable
}

type comparisonDefinition struct {
	x           any
	y           any
	ignoreOrder bool
	exclude     []string
}

// NewScriptComparator creates a new response comparator from the given JavaScript code.
// The script must have a function called 'compare' that accepts two responses and returns a map of diffs.
// The map of diffs can contain anything the user is interested in comparing.
// They can name keys as they want and use the 'diff' function to generate the diff.
func NewScriptComparator(script string) (ScriptComparator, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())

	_, err := vm.RunString(script)
	if err != nil {
		return ScriptComparator{}, errors.New("failed to run the comparator script: " + err.Error())
	}

	compare, ok := goja.AssertFunction(vm.Get("compare"))
	if !ok {
		return ScriptComparator{}, errors.New("compare function not found")
	}

	return ScriptComparator{vm, compare}, nil
}

func (c ScriptComparator) Compare(x, y sender.Response) (map[string][]diffmatchpatch.Diff, error) {
	result, err := c.compare(goja.Undefined(), c.runtime.ToValue(x), c.runtime.ToValue(y))
	if err != nil {
		return nil, errors.New("JavaScript runtime error: " + err.Error())
	}

	compDefs := readComparisonDefinitions(c.runtime, result)
	_ = compDefs

	diffs := make(map[string][]diffmatchpatch.Diff)
	for k, v := range compDefs {
		if d := jsDiff(v.x, v.y, v.ignoreOrder, v.exclude); len(d) != 0 {
			diffs[k] = d
		}
	}

	return diffs, nil
}

func jsDiff(x, y any, ignoreOrder bool, exclude []string) []diffmatchpatch.Diff {
	if x != nil && y != nil && reflect.TypeOf(x).Kind() == reflect.String && reflect.TypeOf(y).Kind() == reflect.String {
		json1 := jsonutils.ReformatJson(x.(string), ignoreOrder, exclude)
		json2 := jsonutils.ReformatJson(y.(string), ignoreOrder, exclude)
		if json1 != json2 {
			return diff.Diff(json1, json2)
		}
		return nil
	}
	json1 := jsonutils.ToJson(x, ignoreOrder, exclude)
	json2 := jsonutils.ToJson(y, ignoreOrder, exclude)
	if json1 != json2 {
		return diff.Diff(json1, json2)
	}
	return nil
}

func readComparisonDefinitions(vm *goja.Runtime, v goja.Value) map[string]comparisonDefinition {
	if v == nil || goja.IsNull(v) || goja.IsUndefined(v) {
		return nil
	}

	obj := v.ToObject(vm)

	defs := make(map[string]comparisonDefinition)
	for _, k := range obj.Keys() {
		value := obj.Get(k)
		def := value.Export()
		if reflect.TypeOf(def).Kind() == reflect.Map {
			defMap := def.(map[string]interface{})
			defs[k] = comparisonDefinition{
				x:           defMap["x"],
				y:           defMap["y"],
				ignoreOrder: getValue[bool](defMap, "ignoreOrder"),
				exclude:     getValues[string](defMap, "exclude"),
			}
		}
	}

	return defs
}

func getValue[T any](m map[string]interface{}, key string) T {
	var zero T
	if v, ok := m[key]; ok && reflect.TypeOf(v).Kind() == reflect.Bool {
		return v.(T)
	}
	return zero
}

func getValues[T any](m map[string]interface{}, key string) []T {
	var result []T
	if v, ok := m[key]; ok && reflect.TypeOf(v).Kind() == reflect.Slice {
		for _, s := range v.([]interface{}) {
			if reflect.TypeOf(s).Kind() == reflect.TypeOf(result).Elem().Kind() {
				result = append(result, s.(T))
			}
		}
	}
	return result
}
