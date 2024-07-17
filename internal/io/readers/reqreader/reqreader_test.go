package reqreader_test

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"testpoint/internal/io/readers/reqreader"
	testutils "testpoint/internal/utils/testing"
)

func TestReadRequestsWithHeader(t *testing.T) {
	tempDir := t.TempDir()
	filename := testutils.CreateTempFile(tempDir, "requests.csv", `
url,method,headers,body
/api/test?prefix=te,PUT,"{""myHeader"":""test1""}","{""field"":""test1""}"
/api/test?prefix=ca,GET,"{""myHeader"":""test2""}","{""field"":""test2""}"
/api/test?prefix=do,DELETE,"{""myHeader"":""test3""}","{""field"":""test3""}"
/api/test?prefix=sp,HEAD,"{""myHeader"":""test4""}","{""field"":""test4""}"
`)

	records := reqreader.ReadRequests(filename, true, 0)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 4 {
		t.Error("incorrect result: expected number of records is 4, got", len(actual))
	}

	expected := []reqreader.ReqRecord{
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=te", "PUT", `{"myHeader":"test1"}`, `{"field":"test1"}`},
			Hash:   11646798338009983096,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=ca", "GET", `{"myHeader":"test2"}`, `{"field":"test2"}`},
			Hash:   16675614452030066654,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=do", "DELETE", `{"myHeader":"test3"}`, `{"field":"test3"}`},
			Hash:   1251291885336478464,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=sp", "HEAD", `{"myHeader":"test4"}`, `{"field":"test4"}`},
			Hash:   17736285143750975039,
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestReadRequestsWithoutHeader(t *testing.T) {
	tempDir := t.TempDir()
	filename := testutils.CreateTempFile(tempDir, "requests.csv", `
/api/test?prefix=te,PUT,"{""myHeader"":""test1""}","{""field"":""test1""}"
/api/test?prefix=ca,GET,"{""myHeader"":""test2""}","{""field"":""test2""}"
/api/test?prefix=do,DELETE,"{""myHeader"":""test3""}","{""field"":""test3""}"
/api/test?prefix=sp,HEAD,"{""myHeader"":""test4""}","{""field"":""test4""}"
`)

	records := reqreader.ReadRequests(filename, false, 0)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 4 {
		t.Error("incorrect result: expected number of records is 4, got", len(actual))
	}

	expected := []reqreader.ReqRecord{
		{
			Values: []string{"/api/test?prefix=te", "PUT", `{"myHeader":"test1"}`, `{"field":"test1"}`},
			Hash:   7213111679473322976,
		},
		{
			Values: []string{"/api/test?prefix=ca", "GET", `{"myHeader":"test2"}`, `{"field":"test2"}`},
			Hash:   16768879361472494806,
		},
		{
			Values: []string{"/api/test?prefix=do", "DELETE", `{"myHeader":"test3"}`, `{"field":"test3"}`},
			Hash:   16488959774387529320,
		},
		{
			Values: []string{"/api/test?prefix=sp", "HEAD", `{"myHeader":"test4"}`, `{"field":"test4"}`},
			Hash:   2507087666846395081,
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestReadRequestsWithEmptyPath(t *testing.T) {
	records := reqreader.ReadRequests("", true, 0)
	actual := testutils.ChanToSlice(records)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}
}

func TestReadRequestsWithEmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	filename := testutils.CreateTempFile(tempDir, "requests.csv", ``)

	records := reqreader.ReadRequests(filename, true, 0)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}
}

func TestReadRequestsWithIncorrectRecords(t *testing.T) {
	tempDir := t.TempDir()
	filename := testutils.CreateTempFile(tempDir, "requests.csv", `
url,method,headers,body
/api/test?prefix=te,PUT,"{""myHeader"":""test1""}","{""field"":""test1""}"
/api/test?prefix=ca,GET,"{""myHeader"":""test2""}"
/api/test?prefix=do,DELETE,"{""myHeader"":""test3""}","{""field"":""test3""}"
/api/test?prefix=sp,HEAD,"{""myHeader"":""test4""}","{""field"":""test4""}""
`)

	records := reqreader.ReadRequests(filename, true, 0)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 2 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}

	expected := []reqreader.ReqRecord{
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=te", "PUT", `{"myHeader":"test1"}`, `{"field":"test1"}`},
			Hash:   11646798338009983096,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=do", "DELETE", `{"myHeader":"test3"}`, `{"field":"test3"}`},
			Hash:   1251291885336478464,
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
/api/test?prefix=te,PUT,"{""myHeader"":""test1""}","{""field"":""test1""}"
/api/test?prefix=ca,GET,"{""myHeader"":""test2""}","{""field"":""test2""}"
/api/test?prefix=do,DELETE,"{""myHeader"":""test3""}","{""field"":""test3""}"
/api/test?prefix=sp,HEAD,"{""myHeader"":""test4""}","{""field"":""test4""}"
`)
	testutils.CreateTempFile(tempDir, "requests-2.csv", `
url,method,headers,body
/api/test2?prefix=am,PUT,"{""myHeader"":""test5""}","{""field"":""test5""}"
/api/test2?prefix=in,GET,"{""myHeader"":""test6""}","{""field"":""test6""}"
/api/test2?prefix=co,DELETE,"{""myHeader"":""test7""}","{""field"":""test7""}"
/api/test2?prefix=st,HEAD,"{""myHeader"":""test8""}","{""field"":""test8""}"
`)

	records := reqreader.ReadRequests(tempDir, true, 0)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 8 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}

	expected := []reqreader.ReqRecord{
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=te", "PUT", `{"myHeader":"test1"}`, `{"field":"test1"}`},
			Hash:   11646798338009983096,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=ca", "GET", `{"myHeader":"test2"}`, `{"field":"test2"}`},
			Hash:   16675614452030066654,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=do", "DELETE", `{"myHeader":"test3"}`, `{"field":"test3"}`},
			Hash:   1251291885336478464,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test?prefix=sp", "HEAD", `{"myHeader":"test4"}`, `{"field":"test4"}`},
			Hash:   17736285143750975039,
		},

		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test2?prefix=am", "PUT", `{"myHeader":"test5"}`, `{"field":"test5"}`},
			Hash:   6261614611056470955,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test2?prefix=in", "GET", `{"myHeader":"test6"}`, `{"field":"test6"}`},
			Hash:   7399124420731000243,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test2?prefix=co", "DELETE", `{"myHeader":"test7"}`, `{"field":"test7"}`},
			Hash:   4441725712074709475,
		},
		{
			Fields: []string{"url", "method", "headers", "body"},
			Values: []string{"/api/test2?prefix=st", "HEAD", `{"myHeader":"test8"}`, `{"field":"test8"}`},
			Hash:   3050802225622638005,
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestReadRequestsFromEmptyDir(t *testing.T) {
	tempDir := t.TempDir()

	records := reqreader.ReadRequests(tempDir, true, 0)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}
}

func TestReadRequestsFromNonexistentDir(t *testing.T) {
	records := reqreader.ReadRequests("/this/directory/does/not/exist/", true, 0)

	actual := testutils.ChanToSlice(records)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}
}
