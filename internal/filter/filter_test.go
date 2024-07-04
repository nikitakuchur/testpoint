package filter_test

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"testpoint/internal/filter"
	"testpoint/internal/reader"
)

func TestFilterWithNoData(t *testing.T) {
	records := make(chan reader.Record)
	close(records)

	filteredRecords := filter.Filter(records)

	var actual = chanToSlice(filteredRecords)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}
}

func TestFilter(t *testing.T) {
	records := make(chan reader.Record)
	go func() {
		records <- reader.Record{Fields: []string{"url", "method"}, Values: []string{"http://test.com/api/test", "GET"}}
		records <- reader.Record{Fields: []string{"url", "method"}, Values: []string{"http://test.com/api/test", "PUT"}}
		records <- reader.Record{Fields: []string{"url", "method"}, Values: []string{"http://foo.com/api/foo", "GET"}}
		records <- reader.Record{Fields: []string{"url", "method"}, Values: []string{"http://bar.com/api/bar", "GET"}}
		records <- reader.Record{Fields: []string{"url", "method"}, Values: []string{"http://test.com/api/test", "GET"}}
		records <- reader.Record{Fields: []string{"url", "method"}, Values: []string{"http://test.com/api/test", "GET"}}
		records <- reader.Record{Fields: []string{"url", "method"}, Values: []string{"http://bar.com/api/bar", "GET"}}
		records <- reader.Record{Fields: []string{"url", "method"}, Values: []string{"http://bar.com/api/bar", "PUT"}}
		close(records)
	}()

	filteredRecords := filter.Filter(records)

	var actual = chanToSlice(filteredRecords)
	if len(actual) != 5 {
		t.Error("incorrect result: expected number of records is 5, got", len(actual))
	}

	expected := []reader.Record{
		{Fields: []string{"url", "method"}, Values: []string{"http://test.com/api/test", "GET"}},
		{Fields: []string{"url", "method"}, Values: []string{"http://test.com/api/test", "PUT"}},
		{Fields: []string{"url", "method"}, Values: []string{"http://foo.com/api/foo", "GET"}},
		{Fields: []string{"url", "method"}, Values: []string{"http://bar.com/api/bar", "GET"}},
		{Fields: []string{"url", "method"}, Values: []string{"http://bar.com/api/bar", "PUT"}},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func chanToSlice(input <-chan reader.Record) []reader.Record {
	var slice []reader.Record
	for rec := range input {
		slice = append(slice, rec)
	}
	return slice
}
