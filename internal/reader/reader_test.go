package reader_test

import (
	"github.com/google/go-cmp/cmp"
	"log"
	"os"
	"testing"
	"testpoint/internal/reader"
)

func TestReadRequestsWithHeader(t *testing.T) {
	tempDir := t.TempDir()
	filename := createTempFile(tempDir, "requests.csv", `
url,method,headers,body
/api/test?prefix=te,PUT,"{""my_header"":""test""}","{""field"":""test""}"
/api/test?prefix=ca,GET,"{""my_header"":""test""}","{""field"":""test""}"
/api/test?prefix=do,GET,"{""my_header"":""test""}","{""field"":""test""}"
/api/test?prefix=sp,GET,"{""my_header"":""test""}","{""field"":""test""}"
`)

	records := reader.ReadRequests(filename, true)

	actual := chanToSlice(records)
	if len(actual) != 4 {
		t.Error("incorrect result: expected slice size is 4, got", len(actual))
	}

	expected := []reader.Record{
		{[]string{"url", "method", "headers", "body"}, []string{"/api/test?prefix=te", "PUT", `{"my_header":"test"}`, `{"field":"test"}`}},
		{[]string{"url", "method", "headers", "body"}, []string{"/api/test?prefix=ca", "GET", `{"my_header":"test"}`, `{"field":"test"}`}},
		{[]string{"url", "method", "headers", "body"}, []string{"/api/test?prefix=do", "GET", `{"my_header":"test"}`, `{"field":"test"}`}},
		{[]string{"url", "method", "headers", "body"}, []string{"/api/test?prefix=sp", "GET", `{"my_header":"test"}`, `{"field":"test"}`}},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestReadRequestsWithoutHeader(t *testing.T) {
	tempDir := t.TempDir()
	filename := createTempFile(tempDir, "requests.csv", `
/api/test?prefix=te,PUT,"{""my_header"":""test""}","{""field"":""test""}"
/api/test?prefix=ca,GET,"{""my_header"":""test""}","{""field"":""test""}"
/api/test?prefix=do,GET,"{""my_header"":""test""}","{""field"":""test""}"
/api/test?prefix=sp,GET,"{""my_header"":""test""}","{""field"":""test""}"
`)

	records := reader.ReadRequests(filename, false)

	actual := chanToSlice(records)
	if len(actual) != 4 {
		t.Error("incorrect result: expected slice size is 4, got", len(actual))
	}

	expected := []reader.Record{
		{Values: []string{"/api/test?prefix=te", "PUT", `{"my_header":"test"}`, `{"field":"test"}`}},
		{Values: []string{"/api/test?prefix=ca", "GET", `{"my_header":"test"}`, `{"field":"test"}`}},
		{Values: []string{"/api/test?prefix=do", "GET", `{"my_header":"test"}`, `{"field":"test"}`}},
		{Values: []string{"/api/test?prefix=sp", "GET", `{"my_header":"test"}`, `{"field":"test"}`}},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestReadRequestsWithEmptyPath(t *testing.T) {
	records := reader.ReadRequests("", true)
	actual := chanToSlice(records)
	if len(actual) != 0 {
		t.Error("incorrect result: expected slice size is 0, got", len(actual))
	}
}

func TestReadRequestsWithEmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	filename := createTempFile(tempDir, "requests.csv", ``)

	records := reader.ReadRequests(filename, true)

	actual := chanToSlice(records)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}
}

func TestReadRequestsWithWrongNumberOfFields(t *testing.T) {
	tempDir := t.TempDir()
	filename := createTempFile(tempDir, "requests.csv", `
url,method,headers,body,test
/api/test?prefix=te,PUT,"{""my_header"":""test""}","{""field"":""test""}"
/api/test?prefix=ca,GET,"{""my_header"":""test""}","{""field"":""test""}"
/api/test?prefix=do,GET,"{""my_header"":""test""}","{""field"":""test""}"
/api/test?prefix=sp,GET,"{""my_header"":""test""}","{""field"":""test""}"
`)

	records := reader.ReadRequests(filename, true)

	actual := chanToSlice(records)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}
}

func TestReadRequestsFromDir(t *testing.T) {
	tempDir := t.TempDir()
	createTempFile(tempDir, "requests-1.csv", `
url,method,headers,body
/api/test?prefix=te,PUT,"{""my_header"":""test""}","{""field"":""test""}"
/api/test?prefix=ca,GET,"{""my_header"":""test""}","{""field"":""test""}"
/api/test?prefix=do,GET,"{""my_header"":""test""}","{""field"":""test""}"
/api/test?prefix=sp,GET,"{""my_header"":""test""}","{""field"":""test""}"
`)
	createTempFile(tempDir, "requests-2.csv", `
url,method,headers,body
/api/test2?prefix=am,PUT,"{""my_header"":""test""}","{""field"":""test""}"
/api/test2?prefix=in,GET,"{""my_header"":""test""}","{""field"":""test""}"
/api/test2?prefix=co,GET,"{""my_header"":""test""}","{""field"":""test""}"
/api/test2?prefix=st,GET,"{""my_header"":""test""}","{""field"":""test""}"
`)

	records := reader.ReadRequests(tempDir, true)

	actual := chanToSlice(records)
	if len(actual) != 8 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}

	expected := []reader.Record{
		{[]string{"url", "method", "headers", "body"}, []string{"/api/test?prefix=te", "PUT", `{"my_header":"test"}`, `{"field":"test"}`}},
		{[]string{"url", "method", "headers", "body"}, []string{"/api/test?prefix=ca", "GET", `{"my_header":"test"}`, `{"field":"test"}`}},
		{[]string{"url", "method", "headers", "body"}, []string{"/api/test?prefix=do", "GET", `{"my_header":"test"}`, `{"field":"test"}`}},
		{[]string{"url", "method", "headers", "body"}, []string{"/api/test?prefix=sp", "GET", `{"my_header":"test"}`, `{"field":"test"}`}},

		{[]string{"url", "method", "headers", "body"}, []string{"/api/test2?prefix=am", "PUT", `{"my_header":"test"}`, `{"field":"test"}`}},
		{[]string{"url", "method", "headers", "body"}, []string{"/api/test2?prefix=in", "GET", `{"my_header":"test"}`, `{"field":"test"}`}},
		{[]string{"url", "method", "headers", "body"}, []string{"/api/test2?prefix=co", "GET", `{"my_header":"test"}`, `{"field":"test"}`}},
		{[]string{"url", "method", "headers", "body"}, []string{"/api/test2?prefix=st", "GET", `{"my_header":"test"}`, `{"field":"test"}`}},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestReadRequestsFromEmptyDir(t *testing.T) {
	tempDir := t.TempDir()

	records := reader.ReadRequests(tempDir, true)

	actual := chanToSlice(records)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}
}

func TestReadRequestsFromNonexistentDir(t *testing.T) {
	records := reader.ReadRequests("/this/directory/does/not/exist/", true)

	actual := chanToSlice(records)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}
}

func createTempFile(dir string, name string, content string) string {
	file, err := os.CreateTemp(dir, name)
	if err != nil {
		log.Fatalln("cannot create a temp file")
	}
	_, err = file.WriteString(content)
	if err != nil {
		log.Fatalln("cannot write into a temp file")
	}
	err = file.Close()
	if err != nil {
		log.Fatalln("cannot close a temp file")
	}
	return file.Name()
}

func chanToSlice(input <-chan reader.Record) []reader.Record {
	var slice []reader.Record
	for rec := range input {
		slice = append(slice, rec)
	}
	return slice
}
