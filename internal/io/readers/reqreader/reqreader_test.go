package reqreader_test

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"testpoint/internal/io/readers/reqreader"
	"testpoint/internal/testutils"
)

func TestReadRequestsWithHeader(t *testing.T) {
	tempDir := t.TempDir()
	filename := testutils.CreateTempFile(tempDir, "requests.csv", `
url,method,headers,body
/api/test?prefix=te,PUT,"{""my_header"":""test1""}","{""field"":""test1""}"
/api/test?prefix=ca,GET,"{""my_header"":""test2""}","{""field"":""test2""}"
/api/test?prefix=do,DELETE,"{""my_header"":""test3""}","{""field"":""test3""}"
/api/test?prefix=sp,HEAD,"{""my_header"":""test4""}","{""field"":""test4""}"
`)

	records := reqreader.ReadRequests(filename, true)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 4 {
		t.Error("incorrect result: expected number of records is 4, got", len(actual))
	}

	expected := []reqreader.ReqRecord{
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=te", "PUT", `{"my_header":"test1"}`, `{"field":"test1"}`},
			Hash:   1267028683842549269,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=ca", "GET", `{"my_header":"test2"}`, `{"field":"test2"}`},
			Hash:   11189344092907974915,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=do", "DELETE", `{"my_header":"test3"}`, `{"field":"test3"}`},
			Hash:   6990969903756593409,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=sp", "HEAD", `{"my_header":"test4"}`, `{"field":"test4"}`},
			Hash:   7828190264647054928,
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestReadRequestsWithoutHeader(t *testing.T) {
	tempDir := t.TempDir()
	filename := testutils.CreateTempFile(tempDir, "requests.csv", `
/api/test?prefix=te,PUT,"{""my_header"":""test1""}","{""field"":""test1""}"
/api/test?prefix=ca,GET,"{""my_header"":""test2""}","{""field"":""test2""}"
/api/test?prefix=do,DELETE,"{""my_header"":""test3""}","{""field"":""test3""}"
/api/test?prefix=sp,HEAD,"{""my_header"":""test4""}","{""field"":""test4""}"
`)

	records := reqreader.ReadRequests(filename, false)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 4 {
		t.Error("incorrect result: expected number of records is 4, got", len(actual))
	}

	expected := []reqreader.ReqRecord{
		{
			Values: []string{"/api/test?prefix=te", "PUT", `{"my_header":"test1"}`, `{"field":"test1"}`},
			Hash:   14472766009977754201,
		},
		{
			Values: []string{"/api/test?prefix=ca", "GET", `{"my_header":"test2"}`, `{"field":"test2"}`},
			Hash:   8160148030485871511,
		},
		{
			Values: []string{"/api/test?prefix=do", "DELETE", `{"my_header":"test3"}`, `{"field":"test3"}`},
			Hash:   14110173358773871789,
		},
		{
			Values: []string{"/api/test?prefix=sp", "HEAD", `{"my_header":"test4"}`, `{"field":"test4"}`},
			Hash:   17692180934814459602,
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestReadRequestsWithEmptyPath(t *testing.T) {
	records := reqreader.ReadRequests("", true)
	actual := testutils.ChanToSlice(records)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}
}

func TestReadRequestsWithEmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	filename := testutils.CreateTempFile(tempDir, "requests.csv", ``)

	records := reqreader.ReadRequests(filename, true)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}
}

func TestReadRequestsWithIncorrectRecords(t *testing.T) {
	tempDir := t.TempDir()
	filename := testutils.CreateTempFile(tempDir, "requests.csv", `
url,method,headers,body
/api/test?prefix=te,PUT,"{""my_header"":""test1""}","{""field"":""test1""}"
/api/test?prefix=ca,GET,"{""my_header"":""test2""}"
/api/test?prefix=do,DELETE,"{""my_header"":""test3""}","{""field"":""test3""}"
/api/test?prefix=sp,HEAD,"{""my_header"":""test4""}","{""field"":""test4""}""
`)

	records := reqreader.ReadRequests(filename, true)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 2 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}

	expected := []reqreader.ReqRecord{
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=te", "PUT", `{"my_header":"test1"}`, `{"field":"test1"}`},
			Hash:   1267028683842549269,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=do", "DELETE", `{"my_header":"test3"}`, `{"field":"test3"}`},
			Hash:   6990969903756593409,
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestReadRequestsFromDir(t *testing.T) {
	tempDir := t.TempDir()
	testutils.CreateTempFile(tempDir, "requests-1.csv", `
url,method,headers,body
/api/test?prefix=te,PUT,"{""my_header"":""test1""}","{""field"":""test1""}"
/api/test?prefix=ca,GET,"{""my_header"":""test2""}","{""field"":""test2""}"
/api/test?prefix=do,DELETE,"{""my_header"":""test3""}","{""field"":""test3""}"
/api/test?prefix=sp,HEAD,"{""my_header"":""test4""}","{""field"":""test4""}"
`)
	testutils.CreateTempFile(tempDir, "requests-2.csv", `
url,method,headers,body
/api/test2?prefix=am,PUT,"{""my_header"":""test5""}","{""field"":""test5""}"
/api/test2?prefix=in,GET,"{""my_header"":""test6""}","{""field"":""test6""}"
/api/test2?prefix=co,DELETE,"{""my_header"":""test7""}","{""field"":""test7""}"
/api/test2?prefix=st,HEAD,"{""my_header"":""test8""}","{""field"":""test8""}"
`)

	records := reqreader.ReadRequests(tempDir, true)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 8 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}

	expected := []reqreader.ReqRecord{
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=te", "PUT", `{"my_header":"test1"}`, `{"field":"test1"}`},
			Hash:   1267028683842549269,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=ca", "GET", `{"my_header":"test2"}`, `{"field":"test2"}`},
			Hash:   11189344092907974915,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=do", "DELETE", `{"my_header":"test3"}`, `{"field":"test3"}`},
			Hash:   6990969903756593409,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=sp", "HEAD", `{"my_header":"test4"}`, `{"field":"test4"}`},
			Hash:   7828190264647054928,
		},

		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test2?prefix=am", "PUT", `{"my_header":"test5"}`, `{"field":"test5"}`},
			Hash:   14942076745510049824,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test2?prefix=in", "GET", `{"my_header":"test6"}`, `{"field":"test6"}`},
			Hash:   10691917426729586672,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test2?prefix=co", "DELETE", `{"my_header":"test7"}`, `{"field":"test7"}`},
			Hash:   5663811562259572916,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test2?prefix=st", "HEAD", `{"my_header":"test8"}`, `{"field":"test8"}`},
			Hash:   17919980904234612742,
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestReadRequestsFromEmptyDir(t *testing.T) {
	tempDir := t.TempDir()

	records := reqreader.ReadRequests(tempDir, true)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}
}

func TestReadRequestsFromNonexistentDir(t *testing.T) {
	records := reqreader.ReadRequests("/this/directory/does/not/exist/", true)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}
}
