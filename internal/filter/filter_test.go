package filter_test

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"testpoint/internal/filter"
	"testpoint/internal/reqreader"
)

func TestFilterWithNoData(t *testing.T) {
	records := make(chan reqreader.Record)
	close(records)

	filteredRecords := filter.Filter(records)

	var actual = chanToSlice(filteredRecords)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}
}

func TestFilter(t *testing.T) {
	records := make(chan reqreader.Record)
	go func() {
		records <- reqreader.Record{Fields: []string{"url", "method"}, Values: []string{"http://test.com/api/test", "GET"}}
		records <- reqreader.Record{Fields: []string{"url", "method"}, Values: []string{"http://test.com/api/test", "PUT"}}
		records <- reqreader.Record{Fields: []string{"url", "method"}, Values: []string{"http://foo.com/api/foo", "GET"}}
		records <- reqreader.Record{Fields: []string{"url", "method"}, Values: []string{"http://bar.com/api/bar", "GET"}}
		records <- reqreader.Record{Fields: []string{"url", "method"}, Values: []string{"http://test.com/api/test", "GET"}}
		records <- reqreader.Record{Fields: []string{"url", "method"}, Values: []string{"http://test.com/api/test", "GET"}}
		records <- reqreader.Record{Fields: []string{"url", "method"}, Values: []string{"http://bar.com/api/bar", "GET"}}
		records <- reqreader.Record{Fields: []string{"url", "method"}, Values: []string{"http://bar.com/api/bar", "PUT"}}
		close(records)
	}()

	filteredRecords := filter.Filter(records)

	var actual = chanToSlice(filteredRecords)
	if len(actual) != 5 {
		t.Error("incorrect result: expected number of records is 5, got", len(actual))
	}

	expected := []reqreader.Record{
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

func chanToSlice(input <-chan reqreader.Record) []reqreader.Record {
	var slice []reqreader.Record
	for rec := range input {
		slice = append(slice, rec)
	}
	return slice
}
