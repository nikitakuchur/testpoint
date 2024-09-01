package strdiff_test

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/nikitakuchur/testpoint/internal/strdiff"
	"testing"
)

func TestCalculateLineDiff(t *testing.T) {
	data := []struct {
		name     string
		text1    string
		text2    string
		expected []strdiff.Diff
	}{
		{
			name:     "empty_strings",
			text1:    "",
			text2:    "",
			expected: nil,
		},
		{
			name:  "one_line_delete",
			text1: "abc",
			text2: "",
			expected: []strdiff.Diff{
				{strdiff.DiffDelete, "abc"},
			},
		},
		{
			name:  "one_line_insert",
			text1: "",
			text2: "abc",
			expected: []strdiff.Diff{
				{strdiff.DiffInsert, "abc"},
			},
		},
		{
			name:  "one_line_equal",
			text1: "abc",
			text2: "abc",
			expected: []strdiff.Diff{
				{strdiff.DiffEqual, "abc"},
			},
		},
		{
			name:  "remove",
			text1: "abc\ndef\nqwerty",
			text2: "abc\nqwerty",
			expected: []strdiff.Diff{
				{strdiff.DiffEqual, "abc\n"},
				{strdiff.DiffDelete, "def\n"},
				{strdiff.DiffEqual, "qwerty"},
			},
		},
		{
			name:  "insert",
			text1: "abc\ndef\nqwerty",
			text2: "abc\ndef\n123\nqwerty",
			expected: []strdiff.Diff{
				{strdiff.DiffEqual, "abc\ndef\n"},
				{strdiff.DiffInsert, "123\n"},
				{strdiff.DiffEqual, "qwerty"},
			},
		},
		{
			name:  "mix",
			text1: "abc\ndef\nqwerty\nhello\nworld",
			text2: "abc\ndef\nhello\ntest\nworld",
			expected: []strdiff.Diff{
				{strdiff.DiffEqual, "abc\ndef\n"},
				{strdiff.DiffDelete, "qwerty\n"},
				{strdiff.DiffEqual, "hello\n"},
				{strdiff.DiffInsert, "test\n"},
				{strdiff.DiffEqual, "world"},
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			actual := strdiff.CalculateLineDiff(d.text1, d.text2)
			if diff := cmp.Diff(d.expected, actual); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestCalculateCharDiff(t *testing.T) {
	data := []struct {
		text1    string
		text2    string
		expected []strdiff.Diff
	}{
		{"", "", nil},
		{"", "abc", []strdiff.Diff{
			{strdiff.DiffInsert, "abc"},
		}},
		{"abc", "", []strdiff.Diff{
			{strdiff.DiffDelete, "abc"},
		}},
		{"abc", "def", []strdiff.Diff{
			{strdiff.DiffDelete, "abc"},
			{strdiff.DiffInsert, "def"},
		}},
		{"abc", "abc", []strdiff.Diff{
			{strdiff.DiffEqual, "abc"},
		}},
		{"abcd", "abc", []strdiff.Diff{
			{strdiff.DiffEqual, "abc"},
			{strdiff.DiffDelete, "d"},
		}},
		{"abc", "abcd", []strdiff.Diff{
			{strdiff.DiffEqual, "abc"},
			{strdiff.DiffInsert, "d"},
		}},
		{"a1b2c3", "abc", []strdiff.Diff{
			{strdiff.DiffEqual, "a"},
			{strdiff.DiffDelete, "1"},
			{strdiff.DiffEqual, "b"},
			{strdiff.DiffDelete, "2"},
			{strdiff.DiffEqual, "c"},
			{strdiff.DiffDelete, "3"},
		}},
		{"gxtxaybd", "aggtfabc", []strdiff.Diff{
			{strdiff.DiffInsert, "a"},
			{strdiff.DiffEqual, "g"},
			{strdiff.DiffDelete, "x"},
			{strdiff.DiffInsert, "g"},
			{strdiff.DiffEqual, "t"},
			{strdiff.DiffDelete, "x"},
			{strdiff.DiffInsert, "f"},
			{strdiff.DiffEqual, "a"},
			{strdiff.DiffDelete, "y"},
			{strdiff.DiffEqual, "b"},
			{strdiff.DiffDelete, "d"},
			{strdiff.DiffInsert, "c"},
		}},
		{"abcdefg", "bdhfg", []strdiff.Diff{
			{strdiff.DiffDelete, "a"},
			{strdiff.DiffEqual, "b"},
			{strdiff.DiffDelete, "c"},
			{strdiff.DiffEqual, "d"},
			{strdiff.DiffDelete, "e"},
			{strdiff.DiffInsert, "h"},
			{strdiff.DiffEqual, "fg"},
		}},
		{"café", "cafe", []strdiff.Diff{
			{strdiff.DiffEqual, "caf"},
			{strdiff.DiffDelete, "é"},
			{strdiff.DiffInsert, "e"},
		}},
		{"🙂🙃😊", "🙃😊", []strdiff.Diff{
			{strdiff.DiffDelete, "🙂"},
			{strdiff.DiffEqual, "🙃😊"},
		}},
		{"abcдфг", "abc", []strdiff.Diff{
			{strdiff.DiffEqual, "abc"},
			{strdiff.DiffDelete, "дфг"},
		}},
		{"hello, world!", "world", []strdiff.Diff{
			{strdiff.DiffDelete, "hello, "},
			{strdiff.DiffEqual, "world"},
			{strdiff.DiffDelete, "!"},
		}},
	}

	for _, d := range data {
		name := fmt.Sprintf("'%s' and '%s'", d.text1, d.text2)
		t.Run(name, func(t *testing.T) {
			actual := strdiff.CalculateCharDiff([]rune(d.text1), []rune(d.text2))
			if diff := cmp.Diff(d.expected, actual); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestDiffToPrettyText(t *testing.T) {
	diff := []strdiff.Diff{
		{strdiff.DiffDelete, "delete"},
		{strdiff.DiffInsert, "insert"},
		{strdiff.DiffEqual, "equal"},
	}

	actual := strdiff.DiffToPrettyText(diff)

	expected := "\x1b[31mdelete\x1b[0m\x1b[32minsert\x1b[0mequal"

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}
