package reporter_test

import (
	"github.com/sergi/go-diff/diffmatchpatch"
	"testing"
	"testpoint/internal/comparator"
	"testpoint/internal/io/readers/respreader"
	"testpoint/internal/io/writers/reporter"
	testutils "testpoint/internal/utils/testing"
)

func TestCsvReporter_Report(t *testing.T) {
	tempDir := t.TempDir()
	rep := reporter.NewCsvReporter(tempDir + "/report.csv")

	diffs := make(chan comparator.RespDiff)

	go func() {
		diffs <- comparator.RespDiff{
			Rec1: respreader.RespRecord{
				ReqUrl: "http://test1.com", ReqMethod: "GET", ReqHeaders: "headers", ReqBody: "body", ReqHash: 123,
				RespStatus: "200", RespBody: "hello",
			},
			Rec2: respreader.RespRecord{
				ReqUrl: "http://test2.com", ReqMethod: "GET", ReqHeaders: "headers", ReqBody: "body", ReqHash: 123,
				RespStatus: "404", RespBody: "not found",
			},
			Diffs: map[string][]diffmatchpatch.Diff{},
		}

		close(diffs)
	}()

	rep.Report(diffs)

	actual := testutils.ReadFile(tempDir + "/report.csv")

	expected := `req1_url,req1_method,req1_headers,req1_body,req2_url,req2_method,req2_headers,req2_body,req_hash,resp1_status,resp1_body,resp2_status,resp2_body
http://test1.com,GET,headers,body,http://test2.com,GET,headers,body,123,200,hello,404,not found
`

	if actual != expected {
		t.Errorf("incorrect result:\nexpected: %v\nactual: %v", expected, actual)
	}
}
