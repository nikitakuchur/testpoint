package comparator

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/nikitakuchur/testpoint/internal/diff"
	"github.com/nikitakuchur/testpoint/internal/sender"
	jsonutils "github.com/nikitakuchur/testpoint/internal/utils/json"
	"github.com/sergi/go-diff/diffmatchpatch"
	"reflect"
	"sync"
)

type ScriptComparator struct {
	runtime     *goja.Runtime
	compare     goja.Callable
	ignoreOrder bool
	mu          sync.Mutex
}

type comparisonDefinition struct {
	x           any
	y           any
	ignoreOrder bool
	exclude     []string
}

// NewScriptComparator creates a new response comparator from the given JavaScript code.
// The script must have a function called 'compare' that accepts two responses and returns a map of comparison definitions.
func NewScriptComparator(script string, ignoreOrder bool) (ScriptComparator, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())

	_, err := vm.RunString(script)
	if err != nil {
		return ScriptComparator{}, fmt.Errorf("failed to run the comparator script: %w", err)
	}

	compare, ok := goja.AssertFunction(vm.Get("compare"))
	if !ok {
		return ScriptComparator{}, errors.New("compare function not found")
	}

	return ScriptComparator{vm, compare, ignoreOrder, sync.Mutex{}}, nil
}

func (c *ScriptComparator) Compare(x, y sender.Response) (map[string][]diffmatchpatch.Diff, error) {
	c.mu.Lock()
	// goja is not thread safe, so we have to lock this piece of code
	result, err := c.compare(goja.Undefined(), c.runtime.ToValue(x), c.runtime.ToValue(y))
	c.mu.Unlock()

	if err != nil {
		return nil, fmt.Errorf("JavaScript runtime error: %w", err)
	}

	compDefs := c.extractComparisonDefinitions(result)

	diffs := make(map[string][]diffmatchpatch.Diff)
	for k, v := range compDefs {
		if d := jsDiff(v.x, v.y, v.ignoreOrder, v.exclude); len(d) != 0 {
			diffs[k] = d
		}
	}

	return diffs, nil
}

func (c *ScriptComparator) extractComparisonDefinitions(v goja.Value) map[string]comparisonDefinition {
	if v == nil || goja.IsNull(v) || goja.IsUndefined(v) {
		return nil
	}

	obj := v.ToObject(c.runtime)

	defs := make(map[string]comparisonDefinition)
	for _, k := range obj.Keys() {
		value := obj.Get(k)
		def := value.Export()
		if reflect.TypeOf(def).Kind() == reflect.Map {
			defMap := def.(map[string]interface{})
			defs[k] = comparisonDefinition{
				x:           defMap["x"],
				y:           defMap["y"],
				ignoreOrder: getValueOrDefault[bool](defMap, "ignoreOrder", c.ignoreOrder),
				exclude:     getValues[string](defMap, "exclude"),
			}
		}
	}

	return defs
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

func getValueOrDefault[T any](m map[string]interface{}, key string, defaultValue T) T {
	if v, ok := m[key]; ok && reflect.TypeOf(v).Kind() == reflect.Bool {
		return v.(T)
	}
	return defaultValue
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
