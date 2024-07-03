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
	scripts := []string{`
function transform(host, record) {
	return {
		url: null,
		method: null,
		headers: null,
		body: null
	};
}
`, `
function transform(host, record) {
	return {
		url: undefined,
		method: undefined,
		headers: undefined,
		body: undefined
	};
}
`, `
function transform(host, record) {
	return {};
}
`, `
function transform(host, record) {
	return null;
}
`, `
function transform(host, record) {
	return undefined;
}
`,
	}

	for _, script := range scripts {
		transformation, _ := transformer.NewTransformation(script)
		record := reader.Record{
			Values: []string{"/api/test", "PUT", "Hello world!"},
		}

		actual, _ := transformation("http://test.com", record)

		expected := sender.Request{}
		if actual != expected {
			t.Error("failed script: ", script)
			t.Errorf("incorrect result: expected request is %v, got %v", expected, actual)
		}
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

func TestDefaultTransformationUrlConcat(t *testing.T) {
	data := []struct {
		name string
		host string
		url  string
	}{
		{"concat1", "http://test.com", "/api/new?param=1&param=2"},
		{"concat2", "http://test.com", "api/new?param=1&param=2"},
		{"concat3", "http://test.com/", "api/new?param=1&param=2"},
		{"concat4", "http://test.com/", "/api/new?param=1&param=2"},
		{"concat5", "http://test.com", "https://site.com/api/new?param=1&param=2"},
		{"concat6", "http://test.com/", "https://site.com/api/new?param=1&param=2"},
		{"concat7", "http://test.com", "https://localhost:8080/api/new?param=1&param=2"},
		{"concat8", "http://test.com/", "https://localhost:8080/api/new?param=1&param=2"},
		{"concat7", "http://test.com/api/new", "https://localhost:8080/api/old?param=1&param=2"},
		{"concat8", "http://test.com/api/new/", "https://localhost:8080/api/old?param=1&param=2"},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			actual, _ := transformer.DefaultTransformation(d.host, reader.Record{Values: []string{d.url}})

			expected := sender.Request{Url: "http://test.com/api/new?param=1&param=2"}
			if actual != expected {
				t.Errorf("incorrect result: expected request is %v, got %v", expected, actual)
			}
		})
	}
}

func TestDefaultTransformationIncorrectHost(t *testing.T) {
	_, err := transformer.DefaultTransformation("://test.com", reader.Record{Values: []string{"/api/test"}})
	if err == nil {
		t.Errorf("incorrect result: expected an error")
	}
}

func TestDefaultTransformationIncorrectUrl(t *testing.T) {
	_, err := transformer.DefaultTransformation("http://test.com", reader.Record{Values: []string{":/api/test"}})
	if err == nil {
		t.Errorf("incorrect result: expected an error")
	}
}
