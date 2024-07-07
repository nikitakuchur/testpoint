package comparator

import (
	"encoding/json"
	"errors"
	"github.com/dop251/goja"
	"github.com/google/go-cmp/cmp"
	"log"
	"testpoint/internal/io/readers/respreader"
)

type RespComparator func(a, b respreader.RespRecord) map[string]string

func NewRespComparator(script string) (RespComparator, error) {
	vm := goja.New()

	_, err := vm.RunString(script)
	if err != nil {
		return nil, errors.New("cannot run the comparator script: " + err.Error())
	}

	compare, ok := goja.AssertFunction(vm.Get("compare"))
	if !ok {
		return nil, errors.New("compare function not found")
	}

	err = vm.Set("diff", func(x, y any) string {
		return cmp.Diff(makePretty(x), makePretty(y))
	})
	if err != nil {
		log.Fatalln("cannot set a diff function for js")
	}

	return func(a, b respreader.RespRecord) map[string]string {
		result, err := compare(goja.Undefined(), vm.ToValue(a.RespBody), vm.ToValue(b.RespBody))
		if err != nil {
			log.Fatal(err)
		}

		return readDiffs(vm, result)
	}, nil
}

func readDiffs(vm *goja.Runtime, value goja.Value) map[string]string {
	obj := value.ToObject(vm)

	m := make(map[string]string)
	for _, k := range obj.Keys() {
		v := obj.Get(k).String()
		if v != "" {
			m[k] = obj.Get(k).String()
		}
	}

	return m
}

func DefaultRespComparator(a, b respreader.RespRecord) map[string]string {
	diffs := make(map[string]string)
	if diff := cmp.Diff(makePretty(a.RespBody), makePretty(b.RespBody)); diff != "" {
		diffs["body"] = diff
	}
	return diffs
}

func makePretty(v any) string {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return ""
	}

	return string(b)
}
