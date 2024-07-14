package comparator_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/sergi/go-diff/diffmatchpatch"
	"testing"
	"testpoint/internal/comparator"
	"testpoint/internal/sender"
)

func TestDefaultRespComparatorWithJsonBody(t *testing.T) {
	rec1 := sender.Response{
		Status: "200",
		Body:   `{"testValue1":"foo","testValue2":"test"}`,
	}
	rec2 := sender.Response{
		Status: "200",
		Body:   `{"testValue1":"bar","testValue2":"test"}`,
	}
	comp := comparator.NewDefaultComparator(false)

	actual, _ := comp.Compare(rec1, rec2)

	expected := map[string][]diffmatchpatch.Diff{
		"body": {
			{Type: diffmatchpatch.DiffEqual, Text: "{\n"},
			{Type: diffmatchpatch.DiffDelete, Text: "  \"testValue1\": \"foo\",\n"},
			{Type: diffmatchpatch.DiffInsert, Text: "  \"testValue1\": \"bar\",\n"},
			{Type: diffmatchpatch.DiffEqual, Text: "  \"testValue2\": \"test\"\n}"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestDefaultRespComparatorWithTextBody(t *testing.T) {
	rec1 := sender.Response{
		Status: "200",
		Body:   "foo",
	}
	rec2 := sender.Response{
		Status: "200",
		Body:   "bar",
	}
	comp := comparator.NewDefaultComparator(false)

	actual, _ := comp.Compare(rec1, rec2)

	expected := map[string][]diffmatchpatch.Diff{
		"body": {
			{Type: diffmatchpatch.DiffDelete, Text: "foo"},
			{Type: diffmatchpatch.DiffInsert, Text: "bar"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestDefaultRespComparatorWithDifferentStatuses(t *testing.T) {
	rec1 := sender.Response{
		Status: "200",
		Body:   "foo",
	}
	rec2 := sender.Response{
		Status: "404",
		Body:   "not found",
	}
	comp := comparator.NewDefaultComparator(false)

	actual, _ := comp.Compare(rec1, rec2)

	expected := map[string][]diffmatchpatch.Diff{
		"status": {
			{Type: diffmatchpatch.DiffDelete, Text: "200"},
			{Type: diffmatchpatch.DiffInsert, Text: "404"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestDefaultRespComparatorWithNoDifferences(t *testing.T) {
	rec1 := sender.Response{
		Status: "200",
		Body:   `{"testValue1":"foo","testValue2":"test"}`,
	}
	rec2 := sender.Response{
		Status: "200",
		Body:   `{"testValue1":"foo","testValue2":"test"}`,
	}
	comp := comparator.NewDefaultComparator(false)

	actual, _ := comp.Compare(rec1, rec2)

	if len(actual) != 0 {
		t.Error("incorrect result: expected number of diffs is 0, got", len(actual))
	}
}
