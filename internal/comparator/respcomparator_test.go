package comparator_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/sergi/go-diff/diffmatchpatch"
	"testing"
	"testpoint/internal/comparator"
	"testpoint/internal/sender"
)

func TestNewRespComparatorWithStatus(t *testing.T) {
	comp, _ := comparator.NewRespComparator(`
function compare(resp1, resp2) {
	return {
		"status": diff(resp1.status, resp2.status)
	};
}
`)

	rec1 := sender.Response{
		Status: "200",
		Body:   `{"testValue1":"foo","testValue2":"test"}`,
	}
	rec2 := sender.Response{
		Status: "404",
		Body:   `not found`,
	}

	actual, _ := comp(rec1, rec2)

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

func TestNewRespComparatorWithBody(t *testing.T) {
	comp, _ := comparator.NewRespComparator(`
function compare(resp1, resp2) {
	return {
		"body": diff(resp1.body, resp2.body)
	};
}
`)

	rec1 := sender.Response{
		Status: "200",
		Body:   `{"testValue1":"foo","testValue2":"test"}`,
	}
	rec2 := sender.Response{
		Status: "200",
		Body:   `{"testValue1":"bar","testValue2":"test"}`,
	}

	actual, _ := comp(rec1, rec2)

	expected := map[string][]diffmatchpatch.Diff{
		"body": {
			{Type: diffmatchpatch.DiffEqual, Text: "{\n"},
			{Type: diffmatchpatch.DiffDelete, Text: "    \"testValue1\": \"foo\",\n"},
			{Type: diffmatchpatch.DiffInsert, Text: "    \"testValue1\": \"bar\",\n"},
			{Type: diffmatchpatch.DiffEqual, Text: "    \"testValue2\": \"test\"\n}"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestNewRespComparator(t *testing.T) {
	comp, _ := comparator.NewRespComparator(`
function compare(resp1, resp2) {
	const body1 = JSON.parse(resp1.body);
	const body2 = JSON.parse(resp2.body);
	return {
		"body.testValue1": diff(body1.testValue1, body2.testValue1),
		"body.testValue2": diff(body1.testValue2, body2.testValue2)
	};
}
`)

	rec1 := sender.Response{
		Status: "200",
		Body:   `{"testValue1":{"value1":"123","value2":"456","value3":"789"},"testValue2":[1,2,3]}`,
	}
	rec2 := sender.Response{
		Status: "200",
		Body:   `{"testValue1":{"value2":"456","value1":"123"},"testValue2":[3,2,1]}`,
	}

	actual, _ := comp(rec1, rec2)

	expected := map[string][]diffmatchpatch.Diff{
		"body.testValue1": {
			{Type: diffmatchpatch.DiffEqual, Text: "{\n    \"value1\": \"123\",\n"},
			{Type: diffmatchpatch.DiffDelete, Text: "    \"value2\": \"456\",\n    \"value3\": \"789\"\n"},
			{Type: diffmatchpatch.DiffInsert, Text: "    \"value2\": \"456\"\n"},
			{Type: diffmatchpatch.DiffEqual, Text: "}"},
		},
		"body.testValue2": {
			{Type: diffmatchpatch.DiffEqual, Text: "[\n"},
			{Type: diffmatchpatch.DiffDelete, Text: "    1,\n"},
			{Type: diffmatchpatch.DiffInsert, Text: "    3,\n"},
			{Type: diffmatchpatch.DiffEqual, Text: "    2,\n"},
			{Type: diffmatchpatch.DiffDelete, Text: "    3\n"},
			{Type: diffmatchpatch.DiffInsert, Text: "    1\n"},
			{Type: diffmatchpatch.DiffEqual, Text: "]"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestNewRespComparatorWithDifferentTypes(t *testing.T) {
	comp, _ := comparator.NewRespComparator(`
function compare(resp1, resp2) {
	return {
		"status": diff("123", 123)
	};
}
`)

	actual, _ := comp(sender.Response{}, sender.Response{})

	expected := map[string][]diffmatchpatch.Diff{
		"status": {
			{Type: diffmatchpatch.DiffDelete, Text: `"123"`},
			{Type: diffmatchpatch.DiffInsert, Text: "123"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestNewRespComparatorWithNull(t *testing.T) {
	comp, _ := comparator.NewRespComparator(`
function compare(resp1, resp2) {
	return {
		"test": diff("123", null)
	};
}
`)

	actual, _ := comp(sender.Response{}, sender.Response{})

	expected := map[string][]diffmatchpatch.Diff{
		"test": {
			{Type: diffmatchpatch.DiffDelete, Text: "\"123\""},
			{Type: diffmatchpatch.DiffInsert, Text: "null"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestNewRespComparatorWithUndefined(t *testing.T) {
	comp, _ := comparator.NewRespComparator(`
function compare(resp1, resp2) {
	return {
		"test": diff("123", undefined)
	};
}
`)

	actual, _ := comp(sender.Response{}, sender.Response{})

	expected := map[string][]diffmatchpatch.Diff{
		"test": {
			{Type: diffmatchpatch.DiffDelete, Text: "\"123\""},
			{Type: diffmatchpatch.DiffInsert, Text: "null"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestNewRespComparatorWithObjects(t *testing.T) {
	comp, _ := comparator.NewRespComparator(`
function compare(resp1, resp2) {
	const foo = {
		test: "123"
	};
	const bar = {
		test: "456"
	};
	return {
		"objects": diff(foo, bar)
	};
}
`)

	actual, _ := comp(sender.Response{}, sender.Response{})

	expected := map[string][]diffmatchpatch.Diff{
		"objects": {
			{Type: diffmatchpatch.DiffEqual, Text: "{\n"},
			{Type: diffmatchpatch.DiffDelete, Text: "    \"test\": \"123\"\n"},
			{Type: diffmatchpatch.DiffInsert, Text: "    \"test\": \"456\"\n"},
			{Type: diffmatchpatch.DiffEqual, Text: "}"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestNewRespComparatorWithArrays(t *testing.T) {
	comp, _ := comparator.NewRespComparator(`
function compare(resp1, resp2) {
	const foo = [1, 2, 3, 4, 5, 6]
	const bar = [1, 0, 3, 4, 5, 6]
	return {
		"arrays": diff(foo, bar)
	};
}
`)

	actual, _ := comp(sender.Response{}, sender.Response{})

	expected := map[string][]diffmatchpatch.Diff{
		"arrays": {
			{Type: diffmatchpatch.DiffEqual, Text: "[\n    1,\n"},
			{Type: diffmatchpatch.DiffDelete, Text: "    2,\n"},
			{Type: diffmatchpatch.DiffInsert, Text: "    0,\n"},
			{Type: diffmatchpatch.DiffEqual, Text: "    3,\n    4,\n    5,\n    6\n]"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestNewRespComparatorWithEmptyDiffs(t *testing.T) {
	data := []struct {
		name   string
		script string
	}{
		{
			"empty",
			`
function compare(resp1, resp2) {
	return {};
}
`,
		},
		{
			"null",
			`
function compare(resp1, resp2) {
	return null;
}
`,
		},
		{
			"undefined",
			`
function compare(resp1, resp2) {
	return undefined;
}
`,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			comp, _ := comparator.NewRespComparator(d.script)
			actual, _ := comp(sender.Response{}, sender.Response{})

			if len(actual) != 0 {
				t.Error("incorrect result: expected number of diffs is 0, got", len(actual))
			}
		})
	}
}

func TestNewRespComparatorWithCreationError(t *testing.T) {
	scripts := []string{"-=24wsfs", ""}
	for _, script := range scripts {
		_, err := comparator.NewRespComparator(script)
		if err == nil {
			t.Errorf("incorrect result: expected an error")
		}
	}
}

func TestNewRespComparatorWithRuntimeError(t *testing.T) {
	comp, _ := comparator.NewRespComparator(`
function compare(resp1, resp2) {
	const a = null;
	a.test();
}
`)

	_, err := comp(sender.Response{}, sender.Response{})

	if err == nil {
		t.Errorf("incorrect result: expected an error")
	}
}

func TestDefaultRespComparatorWithJsonBody(t *testing.T) {
	rec1 := sender.Response{
		Status: "200",
		Body:   `{"testValue1":"foo","testValue2":"test"}`,
	}
	rec2 := sender.Response{
		Status: "200",
		Body:   `{"testValue1":"bar","testValue2":"test"}`,
	}

	actual, _ := comparator.DefaultRespComparator(rec1, rec2)

	expected := map[string][]diffmatchpatch.Diff{
		"body": {
			{Type: diffmatchpatch.DiffEqual, Text: "{\n"},
			{Type: diffmatchpatch.DiffDelete, Text: "    \"testValue1\": \"foo\",\n"},
			{Type: diffmatchpatch.DiffInsert, Text: "    \"testValue1\": \"bar\",\n"},
			{Type: diffmatchpatch.DiffEqual, Text: "    \"testValue2\": \"test\"\n}"},
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

	actual, _ := comparator.DefaultRespComparator(rec1, rec2)

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

	actual, _ := comparator.DefaultRespComparator(rec1, rec2)

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

	actual, _ := comparator.DefaultRespComparator(rec1, rec2)

	if len(actual) != 0 {
		t.Error("incorrect result: expected number of diffs is 0, got", len(actual))
	}
}
