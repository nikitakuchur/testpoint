package comparator_test

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/sergi/go-diff/diffmatchpatch"
	"testing"
	"testpoint/internal/comparator"
	"testpoint/internal/io/readers/respreader"
	"testpoint/internal/sender"
	testutils "testpoint/internal/utils/testing"
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

	diffs := comparator.CompareResponses(records1, records2, comparator.NewDefaultComparator(false), 0)

	var actual = testutils.ChanToSlice(diffs)
	if len(actual) != 2 {
		t.Error("incorrect result: expected number of diffs is 2, got", len(actual))
	}

	expected := []comparator.RespDiff{
		{
			Rec1: respreader.RespRecord{ReqHash: 1, RespStatus: "200", RespBody: "foo"},
			Rec2: respreader.RespRecord{ReqHash: 1, RespStatus: "200", RespBody: "bar"},
			Diffs: map[string][]diffmatchpatch.Diff{
				"body": {
					{Type: diffmatchpatch.DiffDelete, Text: "foo"},
					{Type: diffmatchpatch.DiffInsert, Text: "bar"},
				},
			},
		},
		{
			Rec1: respreader.RespRecord{ReqHash: 2, RespStatus: "500", RespBody: "123"},
			Rec2: respreader.RespRecord{ReqHash: 2, RespStatus: "500", RespBody: "456"},
			Diffs: map[string][]diffmatchpatch.Diff{
				"body": {
					{Type: diffmatchpatch.DiffDelete, Text: "123"},
					{Type: diffmatchpatch.DiffInsert, Text: "456"},
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

	diffs := comparator.CompareResponses(records1, records2, comparator.NewDefaultComparator(false), 0)

	var actual = testutils.ChanToSlice(diffs)
	if len(actual) != 1 {
		t.Error("incorrect result: expected number of diffs is 2, got", len(actual))
	}

	expected := []comparator.RespDiff{
		{
			Rec1: respreader.RespRecord{ReqHash: 1, RespStatus: "200", RespBody: "foo"},
			Rec2: respreader.RespRecord{ReqHash: 1, RespStatus: "200", RespBody: "bar"},
			Diffs: map[string][]diffmatchpatch.Diff{
				"body": {
					{Type: diffmatchpatch.DiffDelete, Text: "foo"},
					{Type: diffmatchpatch.DiffInsert, Text: "bar"},
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

func (ErrorRespComparator) Compare(_, _ sender.Response) (map[string][]diffmatchpatch.Diff, error) {
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

	diffs := comparator.CompareResponses(records1, records2, ErrorRespComparator{}, 0)

	var actual = testutils.ChanToSlice(diffs)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of diffs is 0, got", len(actual))
	}
}
