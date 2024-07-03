package transformer_test

import (
	"testing"
	"testpoint/internal/reader"
	"testpoint/internal/sender"
	"testpoint/internal/transformer"
)

func TestNewTransformationWithoutFields(t *testing.T) {
	transformation, _ := transformer.NewTransformation(`
function transform(host, record) {
	return {
		url: host + record[0],
		method: record[1],
		headers: record[2],
		body: record[3]
	};
}
`)

	record := reader.Record{
		Values: []string{"/api/test", "PUT", `{"test_header":"test_value"}`, "Hello world!"},
	}

	actual, _ := transformation("http://test.com", record)

	expected := sender.Request{
		Url:     "http://test.com/api/test",
		Method:  "PUT",
		Headers: `{"test_header":"test_value"}`,
		Body:    "Hello world!",
	}

	if actual != expected {
		t.Errorf("incorrect result: expected request is %v, got %v", expected, actual)
	}
}

func TestNewTransformationWithFields(t *testing.T) {
	transformation, _ := transformer.NewTransformation(`
function transform(host, record) {
	return {
		url: host + record['url'],
		method: record['method'],
		headers: record['headers'],
		body: record['body']
	};
}
`)

	record := reader.Record{
		Fields: []string{"url", "method", "headers", "body"},
		Values: []string{"/api/test", "PUT", `{"test_header":"test_value"}`, "Hello world!"},
	}

	actual, _ := transformation("http://test.com", record)

	expected := sender.Request{
		Url:     "http://test.com/api/test",
		Method:  "PUT",
		Headers: `{"test_header":"test_value"}`,
		Body:    "Hello world!",
	}

	if actual != expected {
		t.Errorf("incorrect result: expected request is %v, got %v", expected, actual)
	}
}

func TestNewTransformationWithHeadersAsObject(t *testing.T) {
	transformation, _ := transformer.NewTransformation(`
function transform(host, record) {
	return {
		url: host + record[0],
		method: record[1],
		headers: {
			"test_header": "test_value"
		},
		body: record[2]
	};
}
`)

	record := reader.Record{
		Values: []string{"/api/test", "PUT", "Hello world!"},
	}

	actual, _ := transformation("http://test.com", record)

	expected := sender.Request{
		Url:     "http://test.com/api/test",
		Method:  "PUT",
		Headers: `{"test_header":"test_value"}`,
		Body:    "Hello world!",
	}

	if actual != expected {
		t.Errorf("incorrect result: expected request is %v, got %v", expected, actual)
	}
}

func TestNewTransformationWithEmptyRequest(t *testing.T) {
	data := []struct {
		name   string
		script string
	}{
		{
			"null_values",
			`
function transform(host, record) {
	return {
		url: null,
		method: null,
		headers: null,
		body: null
	};
}
`,
		},
		{
			"undefined_values",
			`
function transform(host, record) {
	return {
		url: undefined,
		method: undefined,
		headers: undefined,
		body: undefined
	};
}
`,
		},
		{
			"empty_request",
			`
function transform(host, record) {
	return {};
}
`,
		},
		{
			"null_request",
			`
function transform(host, record) {
	return null;
}
`,
		},
		{
			"undefined_request",
			`
function transform(host, record) {
	return undefined;
}
`,
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			transformation, _ := transformer.NewTransformation(d.script)
			record := reader.Record{
				Values: []string{"/api/test", "PUT", "Hello world!"},
			}

			actual, _ := transformation("http://test.com", record)

			expected := sender.Request{}
			if actual != expected {
				t.Error("failed script: ", d.script)
				t.Errorf("incorrect result: expected request is %v, got %v", expected, actual)
			}
		})
	}
}

func TestNewTransformationWithCreationError(t *testing.T) {
	scripts := []string{"-=24wsfs", ""}
	for _, script := range scripts {
		_, err := transformer.NewTransformation(script)
		if err == nil {
			t.Errorf("incorrect result: expected an error")
		}
	}
}

func TestNewTransformationWithRuntimeError(t *testing.T) {
	transformation, _ := transformer.NewTransformation(`
function transform(host, record) {
	const a = null;
	a.test();
}
`)

	record := reader.Record{
		Values: []string{"/api/test", "PUT", "Hello world!"},
	}

	_, err := transformation("http://test.com", record)

	if err == nil {
		t.Errorf("incorrect result: expected an error")
	}
}

func TestNewTransformationWithMarshallingError(t *testing.T) {
	transformation, _ := transformer.NewTransformation(`
function transform(host, record) {
	const obj = {};
	obj.self = obj;
	return {
		url: host + record[0],
		method: record[1],
		headers: obj,
		body: record[2]
	};
}
`)

	record := reader.Record{
		Values: []string{"/api/test", "PUT", "Hello world!"},
	}

	_, err := transformation("http://test.com", record)

	if err == nil {
		t.Errorf("incorrect result: expected an error")
	}
}

func TestDefaultTransformation(t *testing.T) {
	records := []reader.Record{{
		Values: []string{"/api/test", "PUT", `{"test_header":"test_value"}`, "Hello world!"},
	}, {
		Fields: []string{"body", "headers", "method", "url"},
		Values: []string{"Hello world!", `{"test_header":"test_value"}`, "PUT", "/api/test"},
	}}

	for _, record := range records {
		actual, _ := transformer.DefaultTransformation("http://test.com", record)

		expected := sender.Request{
			Url:     "http://test.com/api/test",
			Method:  "PUT",
			Headers: `{"test_header":"test_value"}`,
			Body:    "Hello world!",
		}
		if actual != expected {
			t.Errorf("incorrect result: expected request is %v, got %v", expected, actual)
		}
	}
}

func TestDefaultTransformationWithFields(t *testing.T) {
	record := reader.Record{
		Fields: []string{"body", "headers", "method", "url"},
		Values: []string{"Hello world!", `{"test_header":"test_value"}`, "PUT", "/api/test"},
	}

	actual, _ := transformer.DefaultTransformation("http://test.com", record)

	expected := sender.Request{
		Url:     "http://test.com/api/test",
		Method:  "PUT",
		Headers: `{"test_header":"test_value"}`,
		Body:    "Hello world!",
	}
	if actual != expected {
		t.Errorf("incorrect result: expected request is %v, got %v", expected, actual)
	}
}

func TestDefaultTransformationWithEmptyRecord(t *testing.T) {
	record := reader.Record{}

	actual, _ := transformer.DefaultTransformation("http://test.com", record)

	expected := sender.Request{Url: "http://test.com"}
	if actual != expected {
		t.Errorf("incorrect result: expected request is %v, got %v", expected, actual)
	}
}

func TestDefaultTransformationWithUrlMerging(t *testing.T) {
	data := []struct {
		name    string
		userUrl string
		reqUrl  string
	}{
		{"right_slash", "http://test.com", "/api/new?param=1&param=2"},
		{"left_slash", "http://test.com/", "api/new?param=1&param=2"},
		{"no_slashes", "http://test.com", "api/new?param=1&param=2"},
		{"both_slashes", "http://test.com/", "/api/new?param=1&param=2"},
		{"host", "http://test.com", "https://site.com/api/new?param=1&param=2"},
		{"host_with_slash", "http://test.com/", "https://site.com/api/new?param=1&param=2"},
		{"host_with_port", "http://test.com", "https://localhost:8080/api/new?param=1&param=2"},
		{"host_with_port_and_slash", "http://test.com/", "https://localhost:8080/api/new?param=1&param=2"},
		{"host_with_path", "http://test.com/api/new", "https://site.com/api/old?param=1&param=2"},
		{"host_with_path_and_slash", "http://test.com/api/new/", "https://site.com/api/old?param=1&param=2"},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			actual, _ := transformer.DefaultTransformation(d.userUrl, reader.Record{Values: []string{d.reqUrl}})

			expected := sender.Request{Url: "http://test.com/api/new?param=1&param=2"}
			if actual != expected {
				t.Errorf("incorrect result: expected request is %v, got %v", expected, actual)
			}
		})
	}
}

func TestDefaultTransformationWithIncorrectUrls(t *testing.T) {
	data := []struct {
		name    string
		userUrl string
		reqUrl  string
	}{
		{"incorrect_user_url", "://test.com", "/api/test"},
		{"incorrect_req_url", "http://test.com", ":/api/test"},
		{"missing_scheme", "test.com", "/api/test"},
		{"missing_host", "http://", "/api/test"},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			_, err := transformer.DefaultTransformation(d.userUrl, reader.Record{Values: []string{d.reqUrl}})
			if err == nil {
				t.Errorf("incorrect result: expected an error")
			}
		})
	}
}
