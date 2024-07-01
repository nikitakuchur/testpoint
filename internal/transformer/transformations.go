package transformer

import (
	"encoding/json"
	"log"
	"strings"
	"testpoint/internal/reader"
)

// DefaultTransformation transforms a raw record to an HTTP request.
// If we don't have a header in the CSV file, the transformation expects the data to be in the following order:
// URL (without the host), HTTP method, headers (in JSON format), body.
// If we do have a header, then it will look for these fields: url, method, headers, and body.
func DefaultTransformation(url string, rec reader.Record) Request {
	params := map[string]string{}
	if rec.Fields != nil {
		// we have a header and we can use it
		for i, field := range rec.Fields {
			params[strings.ToLower(field)] = rec.Values[i]
		}
	} else {
		params["url"] = getValue(rec.Values, 0)
		params["method"] = getValue(rec.Values, 1)
		params["headers"] = getValue(rec.Values, 2)
		params["body"] = getValue(rec.Values, 3)
	}

	return Request{
		Url:     url + params["url"],
		Method:  params["method"],
		Headers: parseHeaders(params["headers"]),
		Body:    params["body"],
	}
}

func parseHeaders(str string) map[string]string {
	var headers map[string]string

	if str == "" {
		return headers
	}

	err := json.Unmarshal([]byte(str), &headers)
	if err != nil {
		log.Fatalln("cannot parse headers:", err)
	}

	return headers
}

func getValue(values []string, i int) string {
	if i >= len(values) {
		return ""
	}
	return values[i]
}
