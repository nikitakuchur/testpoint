package respreader_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/nikitakuchur/testpoint/internal/io/readers/respreader"
	testutils "github.com/nikitakuchur/testpoint/internal/utils/testing"
	"testing"
)

func TestReadResponses(t *testing.T) {
	tempDir := t.TempDir()
	filename := testutils.CreateTempFile(tempDir, "responses.csv", `
req_url,req_method,req_headers,req_body,req_hash,resp_status,resp_body
http://localhost:8080/api/test?prefix=te,PUT,"{""myHeader"":""test1""}","{""field"":""test1""}",123,"200","[1,2,3]"
http://localhost:8080/api/test?prefix=ca,GET,"{""myHeader"":""test2""}","{""field"":""test2""}",234,"404","[4,5,6]"
http://localhost:8080/api/test?prefix=do,DELETE,"{""myHeader"":""test3""}","{""field"":""test3""}",345,"500","[7,8,9]"
http://localhost:8080/api/test?prefix=sp,HEAD,"{""myHeader"":""test4""}","{""field"":""test4""}",456,"201","[10,11,12]"
`)

	records := respreader.ReadResponses(filename)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 4 {
		t.Error("incorrect result: expected slice size is 4, got", len(actual))
	}

	expected := []respreader.RespRecord{
		{
			ReqUrl:     "http://localhost:8080/api/test?prefix=te",
			ReqMethod:  "PUT",
			ReqHeaders: `{"myHeader":"test1"}`,
			ReqBody:    `{"field":"test1"}`,
			ReqHash:    123,
			RespStatus: "200",
			RespBody:   "[1,2,3]",
		},
		{
			ReqUrl:     "http://localhost:8080/api/test?prefix=ca",
			ReqMethod:  "GET",
			ReqHeaders: `{"myHeader":"test2"}`,
			ReqBody:    `{"field":"test2"}`,
			ReqHash:    234,
			RespStatus: "404",
			RespBody:   "[4,5,6]",
		},
		{
			ReqUrl:     "http://localhost:8080/api/test?prefix=do",
			ReqMethod:  "DELETE",
			ReqHeaders: `{"myHeader":"test3"}`,
			ReqBody:    `{"field":"test3"}`,
			ReqHash:    345,
			RespStatus: "500",
			RespBody:   "[7,8,9]",
		},
		{
			ReqUrl:     "http://localhost:8080/api/test?prefix=sp",
			ReqMethod:  "HEAD",
			ReqHeaders: `{"myHeader":"test4"}`,
			ReqBody:    `{"field":"test4"}`,
			ReqHash:    456,
			RespStatus: "201",
			RespBody:   "[10,11,12]",
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestReadRequestsWithWithIncorrectRecords(t *testing.T) {
	tempDir := t.TempDir()
	filename := testutils.CreateTempFile(tempDir, "responses.csv", `
req_url,req_method,req_headers,req_body,req_hash,resp_status,resp_body
http://localhost:8080/api/test?prefix=te,PUT,"{""myHeader"":""test1""}","{""field"":""test1""}",test,"200","[1,2,3]"
http://localhost:8080/api/test?prefix=ca,GET,"{""myHeader"":""test2""}","{""field"":""test2""}",234,"404"
http://localhost:8080/api/test?prefix=do,DELETE,"{""myHeader"":""test3""}","{""field"":""test3""}",345,"500","[7,8,9]"
http://localhost:8080/api/test?prefix=sp,HEAD,"{""myHeader"":""test4""}","{""field"":""test4""}",456,"201","[10,11,12]""
`)

	records := respreader.ReadResponses(filename)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 1 {
		t.Error("incorrect result: expected slice size is 1, got", len(actual))
	}

	expected := []respreader.RespRecord{
		{
			ReqUrl:     "http://localhost:8080/api/test?prefix=do",
			ReqMethod:  "DELETE",
			ReqHeaders: `{"myHeader":"test3"}`,
			ReqBody:    `{"field":"test3"}`,
			ReqHash:    345,
			RespStatus: "500",
			RespBody:   "[7,8,9]",
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestReadRequestsWithEmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	filename := testutils.CreateTempFile(tempDir, "responses.csv", "")

	records := respreader.ReadResponses(filename)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 0 {
		t.Error("incorrect result: expected slice size is 0, got", len(actual))
	}
}
