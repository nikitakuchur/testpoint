package reporter

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"testing"
	"testpoint/internal/comparator"
	"testpoint/internal/io/readers/respreader"
)

func TestLogReporter_Report(t *testing.T) {
	buff := bytes.Buffer{}
	logger := log.New(&buff, "", 0)
	rep := NewLogReporter(logger)

	diffs := make(chan comparator.RespDiff)

	go func() {
		diffs <- comparator.RespDiff{
			Rec1: respreader.RespRecord{
				ReqUrl: "http://test1.com", ReqMethod: "GET", ReqHeaders: "headers", ReqBody: "body", ReqHash: 123, RespStatus: "200",
			},
			Rec2: respreader.RespRecord{
				ReqUrl: "http://test2.com", ReqMethod: "GET", ReqHeaders: "headers", ReqBody: "body", ReqHash: 123, RespStatus: "404",
			},
			Diffs: map[string][]diffmatchpatch.Diff{
				"status": {
					{Type: diffmatchpatch.DiffDelete, Text: "20"},
					{Type: diffmatchpatch.DiffInsert, Text: "4"},
					{Type: diffmatchpatch.DiffEqual, Text: "0"},
					{Type: diffmatchpatch.DiffInsert, Text: "4"},
				},
			},
		}

		close(diffs)
	}()

	rep.Report(diffs)

	expected := "MISMATCH:\n" +
		"req1:\n\turl:\thttp://test1.com\n\tmethod:\tGET\n\theaders:\theaders\n\tbody:\tbody\n\n" +
		"req2:\n\turl:\thttp://test2.com\n\tmethod:\tGET\n\theaders:\theaders\n\tbody:\tbody\n\n" +
		"hash: \t123\nstatus:\n\t\u001B[31m20\u001B[0m\u001B[32m4\u001B[0m0\u001B[32m4\u001B[0m\n"

	if diff := cmp.Diff(expected, buff.String()); diff != "" {
		t.Error(diff)
	}
}

func TestLogReporter_shortenDiffs(t *testing.T) {
	diff := []diffmatchpatch.Diff{
		{
			Type: diffmatchpatch.DiffEqual,
			Text: `penetrate
attract
elegant
marathon
rebellion
overlook
sandwich
venture
incredible`,
		},
		{
			Type: diffmatchpatch.DiffInsert,
			Text: "neighbour",
		},
		{
			Type: diffmatchpatch.DiffEqual,
			Text: `minimum
midnight
graphic
perfect
brother
demonstrate
falsify
election
unlawful
profile`,
		},
		{
			Type: diffmatchpatch.DiffDelete,
			Text: "definite",
		},
		{
			Type: diffmatchpatch.DiffEqual,
			Text: `loyalty
default
overview
outside
humanity
frighten
imagine
thought
accompany
acquisition`,
		},
	}

	actual := shortenDiff(diff)

	expected := []diffmatchpatch.Diff{
		{
			Type: diffmatchpatch.DiffEqual,
			Text: `... // 6 identical lines
sandwich
venture
incredible`,
		},
		{
			Type: diffmatchpatch.DiffInsert,
			Text: "neighbour",
		},
		{
			Type: diffmatchpatch.DiffEqual,
			Text: `minimum
midnight
graphic
... // 4 identical lines
election
unlawful
profile`,
		},
		{
			Type: diffmatchpatch.DiffDelete,
			Text: "definite",
		},
		{
			Type: diffmatchpatch.DiffEqual,
			Text: `loyalty
default
overview
... // 7 identical lines`,
		},
	}

	if d := cmp.Diff(expected, actual); d != "" {
		t.Error(d)
	}
}
