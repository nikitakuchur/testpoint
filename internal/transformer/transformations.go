package transformer

import (
	"errors"
	"github.com/dop251/goja"
	"log"
	"net/url"
	"reflect"
	"strings"
	"testpoint/internal/reader"
	"testpoint/internal/sender"
)

// NewTransformation creates a new transformation from the given JavaScript code.
// The script must have a function called transform that accepts a host and a CSV record, and returns an HTTP request.
func NewTransformation(script string) (Transformation, error) {
	vm := goja.New()

	_, err := vm.RunString(script)
	if err != nil {
		return nil, errors.New("cannot run the transformation script: " + err.Error())
	}

	transform, ok := goja.AssertFunction(vm.Get("transform"))
	if !ok {
		return nil, errors.New("transform function not found")
	}

	return func(host string, rec reader.Record) sender.Request {
		params := createNamedParams(rec)

		var jsRec goja.Value
		if len(params) == 0 {
			jsRec = vm.ToValue(rec.Values)
		} else {
			jsRec = vm.ToValue(params)
		}

		result, err := transform(goja.Undefined(), vm.ToValue(host), jsRec)
		if err != nil {
			// We can't really do much with a runtime error, so let's just log it and skip the request
			log.Println("an error occurred while calling the transform function:", err)
			return sender.Request{}
		}

		if isEmptyValue(result) {
			return sender.Request{}
		}

		obj := result.ToObject(vm)
		return sender.Request{
			Url:     readJsString(obj, "url"),
			Method:  readJsString(obj, "method"),
			Headers: readJsHeaders(vm, obj, "headers"),
			Body:    readJsString(obj, "body"),
		}
	}, nil
}

func readJsString(obj *goja.Object, field string) string {
	v := obj.Get(field)
	if isEmptyValue(v) {
		return ""
	}
	return v.String()
}

func readJsHeaders(vm *goja.Runtime, obj *goja.Object, field string) string {
	v := obj.Get(field)
	if isEmptyValue(v) {
		return ""
	}

	if v.ExportType().Kind() == reflect.String {
		return v.String()
	}

	bytes, err := v.ToObject(vm).MarshalJSON()
	if err != nil {
		// There's a small possibility that we might get a runtime error while converting headers to JSON.
		// For example, it might be caused by a circular structure in headers. Let's log such cases and return an empty string.
		log.Println("an error occurred while reading request headers:", err)
		return ""
	}
	return string(bytes)
}

func isEmptyValue(v goja.Value) bool {
	return v == nil || goja.IsNull(v) || goja.IsUndefined(v)
}

// DefaultTransformation transforms a raw record to an HTTP request.
// If we don't have a header in the CSV file, the transformation expects the data to be in the following order:
// URL (without the host), HTTP method, headers (in JSON format), body.
// If we do have a header, then it will look for these fields: url, method, headers, and body.
func DefaultTransformation(host string, rec reader.Record) sender.Request {
	params := createNamedParams(rec)
	if len(params) == 0 {
		params["url"] = getValue(rec.Values, 0)
		params["method"] = getValue(rec.Values, 1)
		params["headers"] = getValue(rec.Values, 2)
		params["body"] = getValue(rec.Values, 3)
	}

	return sender.Request{
		Url:     buildUrl(host, params["url"]),
		Method:  params["method"],
		Headers: params["headers"],
		Body:    params["body"],
	}
}

func buildUrl(h string, u string) string {
	parsedHost, err := url.Parse(h)
	if err != nil {
		log.Println(err)
		return h + u
	}

	parsedUrl, err := url.Parse(u)
	if err != nil {
		log.Println(err)
		return h + u
	}

	parsedUrl.Scheme = parsedHost.Scheme
	parsedUrl.Host = parsedHost.Host

	return parsedUrl.String()
}

func createNamedParams(rec reader.Record) map[string]string {
	params := map[string]string{}
	if rec.Fields != nil {
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
