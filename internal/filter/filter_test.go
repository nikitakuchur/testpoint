package filter_test

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"testpoint/internal/filter"
	"testpoint/internal/io/readers/reqreader"
)

func TestFilterWithNoData(t *testing.T) {
	records := make(chan reqreader.ReqRecord)
	close(records)

	filteredRecords := filter.Filter(records)

	var actual = chanToSlice(filteredRecords)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of records is 0, got", len(actual))
	}
}

func TestFilter(t *testing.T) {
	records := make(chan reqreader.ReqRecord)
	go func() {
		records <- reqreader.ReqRecord{Fields: []string{"url", "method"}, Values: []string{"http://test.com/api/test", "GET"}, Hash: 1}
		records <- reqreader.ReqRecord{Fields: []string{"url", "method"}, Values: []string{"http://test.com/api/test", "PUT"}, Hash: 2}
		records <- reqreader.ReqRecord{Fields: []string{"url", "method"}, Values: []string{"http://foo.com/api/foo", "GET"}, Hash: 3}
		records <- reqreader.ReqRecord{Fields: []string{"url", "method"}, Values: []string{"http://bar.com/api/bar", "GET"}, Hash: 4}
		records <- reqreader.ReqRecord{Fields: []string{"url", "method"}, Values: []string{"http://test.com/api/test", "GET"}, Hash: 1}
		records <- reqreader.ReqRecord{Fields: []string{"url", "method"}, Values: []string{"http://test.com/api/test", "GET"}, Hash: 1}
		records <- reqreader.ReqRecord{Fields: []string{"url", "method"}, Values: []string{"http://bar.com/api/bar", "GET"}, Hash: 4}
		records <- reqreader.ReqRecord{Fields: []string{"url", "method"}, Values: []string{"http://bar.com/api/bar", "PUT"}, Hash: 5}
		close(records)
	}()

	filteredRecords := filter.Filter(records)

	var actual = chanToSlice(filteredRecords)
	if len(actual) != 5 {
		t.Error("incorrect result: expected number of records is 5, got", len(actual))
	}

	expected := []reqreader.ReqRecord{
		{Fields: []string{"url", "method"}, Values: []string{"http://test.com/api/test", "GET"}, Hash: 1},
		{Fields: []string{"url", "method"}, Values: []string{"http://test.com/api/test", "PUT"}, Hash: 2},
		{Fields: []string{"url", "method"}, Values: []string{"http://foo.com/api/foo", "GET"}, Hash: 3},
		{Fields: []string{"url", "method"}, Values: []string{"http://bar.com/api/bar", "GET"}, Hash: 4},
		{Fields: []string{"url", "method"}, Values: []string{"http://bar.com/api/bar", "PUT"}, Hash: 5},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func chanToSlice(input <-chan reqreader.ReqRecord) []reqreader.ReqRecord {
	var slice []reqreader.ReqRecord
	for rec := range input {
		slice = append(slice, rec)
	}
	return slice
}
