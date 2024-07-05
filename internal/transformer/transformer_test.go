package transformer_test

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"testing"
	"testpoint/internal/reqreader"
	"testpoint/internal/sender"
	"testpoint/internal/transformer"
)

func TestTransformRequestsWithNoData(t *testing.T) {
	records := make(chan reqreader.Record)
	close(records)

	requests := transformer.TransformRequests(nil, records, testTransformation)

	var actual = chanToSlice(requests)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of requests is 0, got", len(actual))
	}
}

func TestTransformRequests(t *testing.T) {
	records := make(chan reqreader.Record)
	go func() {
		records <- reqreader.Record{Values: []string{"/api/test1"}}
		records <- reqreader.Record{Values: []string{"/api/test2"}}
		close(records)
	}()

	requests := transformer.TransformRequests([]string{"http://test1.com", "http://test2.com"}, records, testTransformation)

	var actual = chanToSlice(requests)
	if len(actual) != 4 {
		t.Error("incorrect result: expected number of requests is 4, got", len(actual))
	}

	expected := []sender.Request{
		{Url: "http://test1.com/api/test1", Method: "GET", UserUrl: "http://test1.com"},
		{Url: "http://test2.com/api/test1", Method: "GET", UserUrl: "http://test2.com"},
		{Url: "http://test1.com/api/test2", Method: "GET", UserUrl: "http://test1.com"},
		{Url: "http://test2.com/api/test2", Method: "GET", UserUrl: "http://test2.com"},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestTransformRequestsWithIncorrectRecords(t *testing.T) {
	records := make(chan reqreader.Record)
	go func() {
		records <- reqreader.Record{Values: []string{"/api/test1"}}
		records <- reqreader.Record{Values: []string{"/api/test2"}}
		close(records)
	}()

	requests := transformer.TransformRequests([]string{"http://test1.com", "http://test2.com"}, records, errorTransformation)

	var actual = chanToSlice(requests)
	if len(actual) != 0 {
		t.Error("incorrect result: expected number of requests is 0, got", len(actual))
	}
}

func testTransformation(host string, rec reqreader.Record) (sender.Request, error) {
	return sender.Request{Url: host + rec.Values[0]}, nil
}

func errorTransformation(_ string, _ reqreader.Record) (sender.Request, error) {
	return sender.Request{}, errors.New("error")
}

func chanToSlice(input <-chan sender.Request) []sender.Request {
	var slice []sender.Request
	for req := range input {
		slice = append(slice, req)
	}
	return slice
}
