package reporter_test

import (
	"github.com/sergi/go-diff/diffmatchpatch"
	"testing"
	"testpoint/internal/comparator"
	"testpoint/internal/io/readers/respreader"
	"testpoint/internal/io/writers/reporter"
	"testpoint/internal/testutils"
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

	expected := `req_url_1,req_url_2,req_method,req_headers,req_body,req_hash,resp_status_1,resp_body_1,resp_status_2,resp_body_2
http://test1.com,http://test2.com,GET,headers,body,123,200,hello,404,not found
`

	if actual != expected {
		t.Errorf("incorrect result:\nexpected: %v\nactual: %v", expected, actual)
	}
}
