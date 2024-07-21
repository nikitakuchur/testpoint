package reporter_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/nikitakuchur/testpoint/internal/comparator"
	"github.com/nikitakuchur/testpoint/internal/io/readers/respreader"
	"github.com/nikitakuchur/testpoint/internal/io/writers/reporter"
	"sync"
	"testing"
)

type testReporter struct {
	mismatches []comparator.RespDiff
	mut        sync.Mutex
}

func (r *testReporter) Report(input <-chan comparator.RespDiff) {
	for diff := range input {
		r.mut.Lock()
		r.mismatches = append(r.mismatches, diff)
		r.mut.Unlock()
	}
}

func TestGenerateReportWithNoData(t *testing.T) {
	diffs := make(chan comparator.RespDiff)
	close(diffs)

	rep := testReporter{}

	reporter.GenerateReport(diffs, &rep)

	if len(rep.mismatches) != 0 {
		t.Error("incorrect result: expected number of mismatches is 0, got", len(rep.mismatches))
	}
}

func TestGenerateReport(t *testing.T) {
	input := []comparator.RespDiff{
		{
			Rec1: respreader.RespRecord{ReqUrl: "http://test1.com"},
			Rec2: respreader.RespRecord{ReqUrl: "http://test2.com"},
		},
		{
			Rec1: respreader.RespRecord{ReqUrl: "http://test3.com"},
			Rec2: respreader.RespRecord{ReqUrl: "http://test4.com"},
		},
		{
			Rec1: respreader.RespRecord{ReqUrl: "http://test5.com"},
			Rec2: respreader.RespRecord{ReqUrl: "http://test6.com"},
		},
	}

	diffs := make(chan comparator.RespDiff)
	go func() {
		for _, d := range input {
			diffs <- d
		}
		close(diffs)
	}()

	reps := []reporter.Reporter{
		&testReporter{},
		&testReporter{},
		&testReporter{},
	}

	reporter.GenerateReport(diffs, reps...)

	for _, rep := range reps {
		testRep := rep.(*testReporter)
		if d := cmp.Diff(input, testRep.mismatches); d != "" {
			t.Error(d)
		}
	}
}
