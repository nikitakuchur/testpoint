package comparator_test

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/nikitakuchur/testpoint/internal/comparator"
	"github.com/nikitakuchur/testpoint/internal/io/readers/respreader"
	"github.com/nikitakuchur/testpoint/internal/sender"
	"github.com/nikitakuchur/testpoint/internal/strdiff"
	testutils "github.com/nikitakuchur/testpoint/internal/utils/testing"
	"testing"
)

func TestCompareResponses(t *testing.T) {
	records1 := make(chan respreader.RespRecord)
	records2 := make(chan respreader.RespRecord)

	go func() {
		records1 <- respreader.RespRecord{ReqHash: 1, RespStatus: "200", RespBody: "foo"}
		records2 <- respreader.RespRecord{ReqHash: 3, RespStatus: "404", RespBody: "not found"}
		records1 <- respreader.RespRecord{ReqHash: 2, RespStatus: "500", RespBody: "123"}
		records2 <- respreader.RespRecord{ReqHash: 1, RespStatus: "200", RespBody: "bar"}
		records1 <- respreader.RespRecord{ReqHash: 3, RespStatus: "404", RespBody: "not found"}
		records2 <- respreader.RespRecord{ReqHash: 2, RespStatus: "500", RespBody: "456"}
		close(records1)
		close(records2)
	}()

	diffs := comparator.CompareResponses(records1, records2, 0, comparator.NewDefaultComparator(false), 1)

	var actual = testutils.ChanToSlice(diffs)
	if len(actual) != 2 {
		t.Error("incorrect result: expected number of diffs is 2, got", len(actual))
	}

	expected := []comparator.RespDiff{
		{
			Rec1: respreader.RespRecord{ReqHash: 1, RespStatus: "200", RespBody: "foo"},
			Rec2: respreader.RespRecord{ReqHash: 1, RespStatus: "200", RespBody: "bar"},
			Diffs: map[string][]strdiff.Diff{
				"body": {
					{Operation: strdiff.DiffDelete, Text: "foo"},
					{Operation: strdiff.DiffInsert, Text: "bar"},
				},
			},
		},
		{
			Rec1: respreader.RespRecord{ReqHash: 2, RespStatus: "500", RespBody: "123"},
			Rec2: respreader.RespRecord{ReqHash: 2, RespStatus: "500", RespBody: "456"},
			Diffs: map[string][]strdiff.Diff{
				"body": {
					{Operation: strdiff.DiffDelete, Text: "123"},
					{Operation: strdiff.DiffInsert, Text: "456"},
				},
			},
		},
	}

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Error(diff)
	}
}

func TestCompareResponsesWithMissingRecords(t *testing.T) {
	records1 := make(chan respreader.RespRecord)
	records2 := make(chan respreader.RespRecord)

	go func() {
		records1 <- respreader.RespRecord{ReqHash: 1, RespStatus: "200", RespBody: "foo"}
		records1 <- respreader.RespRecord{ReqHash: 2, RespStatus: "500", RespBody: "123"}
		records2 <- respreader.RespRecord{ReqHash: 1, RespStatus: "200", RespBody: "bar"}
		records1 <- respreader.RespRecord{ReqHash: 3, RespStatus: "404", RespBody: "not found"}
		close(records1)
		close(records2)
	}()

	diffs := comparator.CompareResponses(records1, records2, 0, comparator.NewDefaultComparator(false), 1)

	var actual = testutils.ChanToSlice(diffs)
	if len(actual) != 1 {
		t.Error("incorrect result: expected number of diffs is 2, got", len(actual))
	}

	expected := []comparator.RespDiff{
		{
			Rec1: respreader.RespRecord{ReqHash: 1, RespStatus: "200", RespBody: "foo"},
			Rec2: respreader.RespRecord{ReqHash: 1, RespStatus: "200", RespBody: "bar"},
			Diffs: map[string][]strdiff.Diff{
				"body": {
					{Operation: strdiff.DiffDelete, Text: "foo"},
					{Operation: strdiff.DiffInsert, Text: "bar"},
				},
			},
		},
	}

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Error(diff)
	}
}

type ErrorRespComparator struct {
}

func (ErrorRespComparator) Compare(_, _ sender.Response) (map[string][]strdiff.Diff, error) {
	return nil, errors.New("error")
}

func TestCompareResponsesWithErrors(t *testing.T) {
	records1 := make(chan respreader.RespRecord)
	records2 := make(chan respreader.RespRecord)

	go func() {
		records1 <- respreader.RespRecord{ReqHash: 1, RespStatus: "200", RespBody: "foo"}
		records2 <- respreader.RespRecord{ReqHash: 3, RespStatus: "404", RespBody: "not found"}
		records1 <- respreader.RespRecord{ReqHash: 2, RespStatus: "500", RespBody: "123"}
		records2 <- respreader.RespRecord{ReqHash: 1, RespStatus: "200", RespBody: "bar"}
		records1 <- respreader.RespRecord{ReqHash: 3, RespStatus: "404", RespBody: "not found"}
		records2 <- respreader.RespRecord{ReqHash: 2, RespStatus: "500", RespBody: "456"}
		close(records1)
		close(records2)
	}()

	diffs := comparator.CompareResponses(records1, records2, 0, ErrorRespComparator{}, 1)

	var actual = testutils.ChanToSlice(diffs)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of diffs is 0, got", len(actual))
	}
}
