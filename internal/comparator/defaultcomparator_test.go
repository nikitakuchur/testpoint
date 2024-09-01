package comparator_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/nikitakuchur/testpoint/internal/comparator"
	"github.com/nikitakuchur/testpoint/internal/sender"
	"github.com/nikitakuchur/testpoint/internal/strdiff"
	"testing"
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

	expected := map[string][]strdiff.Diff{
		"body": {
			{Operation: strdiff.DiffEqual, Text: "{\n"},
			{Operation: strdiff.DiffDelete, Text: "  \"testValue1\": \"foo\",\n"},
			{Operation: strdiff.DiffInsert, Text: "  \"testValue1\": \"bar\",\n"},
			{Operation: strdiff.DiffEqual, Text: "  \"testValue2\": \"test\"\n}"},
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

	expected := map[string][]strdiff.Diff{
		"body": {
			{Operation: strdiff.DiffDelete, Text: "foo"},
			{Operation: strdiff.DiffInsert, Text: "bar"},
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

	expected := map[string][]strdiff.Diff{
		"status": {
			{Operation: strdiff.DiffDelete, Text: "200"},
			{Operation: strdiff.DiffInsert, Text: "404"},
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
