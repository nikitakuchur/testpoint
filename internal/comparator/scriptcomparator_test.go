package comparator_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/nikitakuchur/testpoint/internal/comparator"
	"github.com/nikitakuchur/testpoint/internal/sender"
	"github.com/nikitakuchur/testpoint/internal/strdiff"
	"testing"
)

func TestNewScriptComparatorWithStatus(t *testing.T) {
	comp, _ := comparator.NewScriptComparator(`
function compare(resp1, resp2) {
	return {
		"status": {x:resp1.status, y:resp2.status}
	};
}
`, false)

	rec1 := sender.Response{
		Status: "200",
		Body:   `{"testValue1":"foo","testValue2":"test"}`,
	}
	rec2 := sender.Response{
		Status: "404",
		Body:   `not found`,
	}

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

func TestNewScriptComparatorWithBody(t *testing.T) {
	comp, _ := comparator.NewScriptComparator(`
function compare(resp1, resp2) {
	return {
		"body": {x: resp1.body, y: resp2.body}
	};
}
`, false)

	rec1 := sender.Response{
		Status: "200",
		Body:   `{"testValue1":"foo","testValue2":"test"}`,
	}
	rec2 := sender.Response{
		Status: "200",
		Body:   `{"testValue1":"bar","testValue2":"test"}`,
	}

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

func TestNewScriptComparator(t *testing.T) {
	comp, _ := comparator.NewScriptComparator(`
function compare(resp1, resp2) {
	const body1 = JSON.parse(resp1.body);
	const body2 = JSON.parse(resp2.body);
	return {
		"body.testValue1": {x: body1.testValue1, y: body2.testValue1},
		"body.testValue2": {x: body1.testValue2, y: body2.testValue2}
	};
}
`, false)

	rec1 := sender.Response{
		Status: "200",
		Body:   `{"testValue1":{"value1":"123","value2":"456","value3":"789"},"testValue2":[1,2,3]}`,
	}
	rec2 := sender.Response{
		Status: "200",
		Body:   `{"testValue1":{"value2":"456","value1":"123"},"testValue2":[3,2,1]}`,
	}

	actual, _ := comp.Compare(rec1, rec2)

	expected := map[string][]strdiff.Diff{
		"body.testValue1": {
			{Operation: strdiff.DiffEqual, Text: "{\n  \"value1\": \"123\",\n"},
			{Operation: strdiff.DiffDelete, Text: "  \"value2\": \"456\",\n  \"value3\": \"789\"\n"},
			{Operation: strdiff.DiffInsert, Text: "  \"value2\": \"456\"\n"},
			{Operation: strdiff.DiffEqual, Text: "}"},
		},
		"body.testValue2": {
			{Operation: strdiff.DiffEqual, Text: "[\n"},
			{Operation: strdiff.DiffDelete, Text: "  1,\n"},
			{Operation: strdiff.DiffInsert, Text: "  3,\n"},
			{Operation: strdiff.DiffEqual, Text: "  2,\n"},
			{Operation: strdiff.DiffDelete, Text: "  3\n"},
			{Operation: strdiff.DiffInsert, Text: "  1\n"},
			{Operation: strdiff.DiffEqual, Text: "]"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestNewScriptComparatorWithEqualValues(t *testing.T) {
	comp, _ := comparator.NewScriptComparator(`
function compare(resp1, resp2) {
	return {
		"status": {x: "123", y: "123"}
	};
}
`, false)

	actual, _ := comp.Compare(sender.Response{}, sender.Response{})

	if len(actual) != 0 {
		t.Error("incorrect result: expected number of diffs is 0, got", len(actual))
	}
}

func TestNewScriptComparatorWithDifferentTypes(t *testing.T) {
	comp, _ := comparator.NewScriptComparator(`
function compare(resp1, resp2) {
	return {
		"status": {x: "123", y: 123}
	};
}
`, false)

	actual, _ := comp.Compare(sender.Response{}, sender.Response{})

	expected := map[string][]strdiff.Diff{
		"status": {
			{Operation: strdiff.DiffDelete, Text: `"123"`},
			{Operation: strdiff.DiffInsert, Text: "123"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestNewScriptComparatorWithNull(t *testing.T) {
	comp, _ := comparator.NewScriptComparator(`
function compare(resp1, resp2) {
	return {
		"test": {x: "123", y: null}
	};
}
`, false)

	actual, _ := comp.Compare(sender.Response{}, sender.Response{})

	expected := map[string][]strdiff.Diff{
		"test": {
			{Operation: strdiff.DiffDelete, Text: "\"123\""},
			{Operation: strdiff.DiffInsert, Text: "null"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestNewScriptComparatorWithUndefined(t *testing.T) {
	comp, _ := comparator.NewScriptComparator(`
function compare(resp1, resp2) {
	return {
		"test": {x: "123", y: undefined}
	};
}
`, false)

	actual, _ := comp.Compare(sender.Response{}, sender.Response{})

	expected := map[string][]strdiff.Diff{
		"test": {
			{Operation: strdiff.DiffDelete, Text: "\"123\""},
			{Operation: strdiff.DiffInsert, Text: "null"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestNewScriptComparatorWithObjects(t *testing.T) {
	comp, _ := comparator.NewScriptComparator(`
function compare(resp1, resp2) {
	const foo = {
		test: "123"
	};
	const bar = {
		test: "456"
	};
	return {
		"objects": {x: foo, y: bar}
	};
}
`, false)

	actual, _ := comp.Compare(sender.Response{}, sender.Response{})

	expected := map[string][]strdiff.Diff{
		"objects": {
			{Operation: strdiff.DiffEqual, Text: "{\n"},
			{Operation: strdiff.DiffDelete, Text: "  \"test\": \"123\"\n"},
			{Operation: strdiff.DiffInsert, Text: "  \"test\": \"456\"\n"},
			{Operation: strdiff.DiffEqual, Text: "}"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestNewScriptComparatorWithArrays(t *testing.T) {
	comp, _ := comparator.NewScriptComparator(`
function compare(resp1, resp2) {
	const foo = [1, 2, 3, 4, 5, 6]
	const bar = [1, 0, 3, 4, 5, 6]
	return {
		"arrays": {x: foo, y: bar}
	};
}
`, false)

	actual, _ := comp.Compare(sender.Response{}, sender.Response{})

	expected := map[string][]strdiff.Diff{
		"arrays": {
			{Operation: strdiff.DiffEqual, Text: "[\n  1,\n"},
			{Operation: strdiff.DiffDelete, Text: "  2,\n"},
			{Operation: strdiff.DiffInsert, Text: "  0,\n"},
			{Operation: strdiff.DiffEqual, Text: "  3,\n  4,\n  5,\n  6\n]"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestNewScriptComparatorWithEmptyDiffs(t *testing.T) {
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
			comp, _ := comparator.NewScriptComparator(d.script, false)
			actual, _ := comp.Compare(sender.Response{}, sender.Response{})

			if len(actual) != 0 {
				t.Error("incorrect result: expected number of diffs is 0, got", len(actual))
			}
		})
	}
}

func TestNewScriptComparatorWithCreationError(t *testing.T) {
	scripts := []string{"-=24wsfs", ""}
	for _, script := range scripts {
		_, err := comparator.NewScriptComparator(script, false)
		if err == nil {
			t.Errorf("incorrect result: expected an error")
		}
	}
}

func TestNewScriptComparatorWithRuntimeError(t *testing.T) {
	comp, _ := comparator.NewScriptComparator(`
function compare(resp1, resp2) {
	const a = null;
	a.test();
}
`, false)

	_, err := comp.Compare(sender.Response{}, sender.Response{})

	if err == nil {
		t.Errorf("incorrect result: expected an error")
	}
}
