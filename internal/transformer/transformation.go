package transformer

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"net/url"
	"reflect"
	"strings"
	"testpoint/internal/io/readers/reqreader"
	"testpoint/internal/sender"
)

type Transformation func(userUrl string, rec reqreader.ReqRecord) (sender.Request, error)

// NewTransformation creates a new transformation from the given JavaScript code.
// The script must have a function called transform that accepts a user url and a CSV record, and returns an HTTP request.
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

	return func(userUrl string, rec reqreader.ReqRecord) (sender.Request, error) {
		params := createNamedParams(rec)

		var jsRec goja.Value
		if len(params) == 0 {
			jsRec = vm.ToValue(rec.Values)
		} else {
			jsRec = vm.ToValue(params)
		}

		result, err := transform(goja.Undefined(), vm.ToValue(userUrl), jsRec)
		if err != nil {
			// We can't really do much with a runtime error, so let's just return an error to skip the record
			return sender.Request{}, errors.New("JavaScript runtime error: " + err.Error())
		}

		if isEmptyValue(result) {
			return sender.Request{}, nil
		}

		obj := result.ToObject(vm)

		parsedHeaders, err := readJsHeaders(vm, obj, "headers")
		if err != nil {
			return sender.Request{}, errors.New("JavaScript runtime error: " + err.Error())
		}

		return sender.Request{
			Url:     readJsString(obj, "url"),
			Method:  readJsString(obj, "method"),
			Headers: parsedHeaders,
			Body:    readJsString(obj, "body"),
		}, nil
	}, nil
}

func readJsString(obj *goja.Object, field string) string {
	v := obj.Get(field)
	if isEmptyValue(v) {
		return ""
	}
	return v.String()
}

func readJsHeaders(vm *goja.Runtime, obj *goja.Object, field string) (string, error) {
	v := obj.Get(field)
	if isEmptyValue(v) {
		return "", nil
	}

	if v.ExportType().Kind() == reflect.String {
		return v.String(), nil
	}

	bytes, err := v.ToObject(vm).MarshalJSON()
	if err != nil {
		// There's a small possibility that we might get a runtime error while converting headers to JSON.
		// For example, it might be caused by marshalling a circular structure.
		return "", err
	}
	return string(bytes), nil
}

func isEmptyValue(v goja.Value) bool {
	return v == nil || goja.IsNull(v) || goja.IsUndefined(v)
}

// DefaultTransformation transforms a raw record to an HTTP request.
// If we don't have a header in the CSV file, the transformation expects the data to be in the following order:
// URL, HTTP method, headers (in JSON format), body.
// If we do have a header, then it will look for these fields: url, method, headers, and body.
func DefaultTransformation(userUrl string, rec reqreader.ReqRecord) (sender.Request, error) {
	params := createNamedParams(rec)
	if len(params) == 0 {
		params["url"] = getValue(rec.Values, 0)
		params["method"] = getValue(rec.Values, 1)
		params["headers"] = getValue(rec.Values, 2)
		params["body"] = getValue(rec.Values, 3)
	}

	mergedUrl, err := mergeUrls(params["url"], userUrl)
	if err != nil {
		// If the URL cannot be parsed, it's better to return an error and skip the record
		return sender.Request{}, err
	}

	return sender.Request{
		Url:     mergedUrl,
		Method:  params["method"],
		Headers: params["headers"],
		Body:    params["body"],
	}, nil
}

// mergeUrls merges request URLs from the input files with the user's URL.
// For example, let's assume we have the following URL in the file: "http://test.com/api/old?param=123".
// If the user's URL is "http://newtest.com", this function will return "http://newtest.com/api/old?param=123".
// If the user's URL is "http://newtest.com/api/new", this function will return "http://newtest.com/api/new?param=123".
func mergeUrls(requestUrl string, userUrl string) (string, error) {
	parsedRequestUrl, err := url.Parse(requestUrl)
	if err != nil {
		return "", err
	}

	if userUrl == "" {
		return requestUrl, nil
	}

	parsedUserUrl, err := url.Parse(userUrl)
	if err != nil {
		return "", err
	}

	if parsedUserUrl.Scheme == "" {
		return "", errors.New(fmt.Sprintf("parse \"%s\": missing protocol scheme", userUrl))
	}
	if parsedUserUrl.Host == "" {
		return "", errors.New(fmt.Sprintf("parse \"%s\": missing host", userUrl))
	}

	parsedRequestUrl.Scheme = parsedUserUrl.Scheme
	parsedRequestUrl.Host = parsedUserUrl.Host
	newPath := strings.TrimRight(parsedUserUrl.Path, "/")
	if newPath != "" {
		parsedRequestUrl.Path = newPath
	}

	return parsedRequestUrl.String(), nil
}

func createNamedParams(rec reqreader.ReqRecord) map[string]string {
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
