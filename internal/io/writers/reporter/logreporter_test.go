package reporter

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"github.com/nikitakuchur/testpoint/internal/comparator"
	"github.com/nikitakuchur/testpoint/internal/io/readers/respreader"
	"github.com/nikitakuchur/testpoint/internal/strdiff"
	"log"
	"testing"
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
			Diffs: map[string][]strdiff.Diff{
				"status": {
					{Operation: strdiff.DiffDelete, Text: "20"},
					{Operation: strdiff.DiffInsert, Text: "4"},
					{Operation: strdiff.DiffEqual, Text: "0"},
					{Operation: strdiff.DiffInsert, Text: "4"},
				},
			},
		}

		close(diffs)
	}()

	rep.Report(diffs)

	expected := "MISMATCH:\n" +
		"req1:\n\turl: http://test1.com\n\tmethod: GET\n\theaders: headers\n\tbody: body\n" +
		"req2:\n\turl: http://test2.com\n\tmethod: GET\n\theaders: headers\n\tbody: body\n\n" +
		"hash: \t123\n\nstatus:\n\t\u001B[31m20\n\t\u001B[0m\u001B[32m4\n\t\u001B[0m0\n\t\u001B[32m4\n\t\u001B[0m\n"

	if diff := cmp.Diff(expected, buff.String()); diff != "" {
		t.Error(diff)
	}
}

func TestLogReporter_shortenDiffs(t *testing.T) {
	diff := []strdiff.Diff{
		{
			Operation: strdiff.DiffEqual,
			Text: `penetrate
attract
elegant
marathon
rebellion
overlook
sandwich
venture
incredible
`,
		},
		{
			Operation: strdiff.DiffInsert,
			Text:      "neighbour",
		},
		{
			Operation: strdiff.DiffEqual,
			Text: `minimum
midnight
graphic
perfect
brother
demonstrate
falsify
election
unlawful
profile
`,
		},
		{
			Operation: strdiff.DiffDelete,
			Text:      "definite",
		},
		{
			Operation: strdiff.DiffEqual,
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

	expected := []strdiff.Diff{
		{
			Operation: strdiff.DiffEqual,
			Text: `... // 6 identical lines
sandwich
venture
incredible
`,
		},
		{
			Operation: strdiff.DiffInsert,
			Text:      "neighbour\n",
		},
		{
			Operation: strdiff.DiffEqual,
			Text: `minimum
midnight
graphic
... // 4 identical lines
election
unlawful
profile
`,
		},
		{
			Operation: strdiff.DiffDelete,
			Text:      "definite\n",
		},
		{
			Operation: strdiff.DiffEqual,
			Text: `loyalty
default
overview
... // 7 identical lines
`,
		},
	}

	if d := cmp.Diff(expected, actual); d != "" {
		t.Error(d)
	}
}
