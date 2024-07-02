package transformer

import (
	"github.com/dop251/goja"
	"log"
	"os"
	"strings"
	"testpoint/internal/reader"
	"testpoint/internal/sender"
)

// NewTransformation creates a new transformation from the given JavaScript file.
// The script must have a function called transform that accepts a host and a CSV record, and returns an HTTP request.
func NewTransformation(filename string) Transformation {
	script := readScript(filename)

	vm := goja.New()
	_, err := vm.RunString(script)
	if err != nil {
		log.Fatalln("cannot run the transformation script:", err)
	}
	transform, ok := goja.AssertFunction(vm.Get("transform"))
	if !ok {
		log.Fatalln("transform function not found!")
	}

	return func(host string, rec reader.Record) sender.Request {
		params := createParams(rec)

		var jsRec goja.Value
		if len(params) == 0 {
			jsRec = vm.ToValue(rec.Values)
		} else {
			jsRec = vm.ToValue(params)
		}
		function, err := transform(goja.Undefined(), vm.ToValue(host), jsRec)
		if err != nil {
			log.Fatalln("cannot call the transform function:", err)
		}

		obj := function.ToObject(vm)
		return sender.Request{
			Url:     readJsValue(obj, "url", ""),
			Method:  readJsValue(obj, "method", "GET"),
			Headers: readJsHeaders(vm, obj, "headers"),
			Body:    readJsValue(obj, "body", ""),
		}
	}
}

func readScript(filename string) string {
	script, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln("cannot read the transformation script:", err)
	}
	return string(script)
}

func readJsValue(obj *goja.Object, field string, def string) string {
	v := obj.Get(field)
	if v == nil {
		return def
	}
	return v.String()
}

func readJsHeaders(vm *goja.Runtime, obj *goja.Object, field string) string {
	v := obj.Get(field)
	if v == nil {
		return ""
	}
	// TODO: doesn't work if the headers are null
	bytes, err := v.ToObject(vm).MarshalJSON()
	if err != nil {
		log.Fatalln("The request headers are not an object:", err)
	}
	return string(bytes)
}

// DefaultTransformation transforms a raw record to an HTTP request.
// If we don't have a header in the CSV file, the transformation expects the data to be in the following order:
// URL (without the host), HTTP method, headers (in JSON format), body.
// If we do have a header, then it will look for these fields: url, method, headers, and body.
func DefaultTransformation(host string, rec reader.Record) sender.Request {
	params := createParams(rec)
	if len(params) == 0 {
		params["url"] = getValue(rec.Values, 0)
		params["method"] = getValue(rec.Values, 1)
		params["headers"] = getValue(rec.Values, 2)
		params["body"] = getValue(rec.Values, 3)
	}

	if params["method"] == "" {
		params["method"] = "GET"
	}

	return sender.Request{
		Url:     host + params["url"],
		Method:  params["method"],
		Headers: params["headers"],
		Body:    params["body"],
	}
}

func createParams(rec reader.Record) map[string]string {
	params := map[string]string{}
	if rec.Fields != nil {
		// we have a header and we can use it
		for i, field := range rec.Fields {
			params[strings.ToLower(field)] = rec.Values[i]
		}
	}
	return params
}

func getValue(values []string, i int) string {
	if i >= len(values) {
		return ""
	}
	return values[i]
}
