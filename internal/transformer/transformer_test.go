package transformer_test

import (
	"errors"
	"testing"
	"testpoint/internal/reader"
	"testpoint/internal/sender"
	"testpoint/internal/transformer"
)

func TestTransformRequestsWithNoData(t *testing.T) {
	records := make(chan reader.Record)
	close(records)

	requests := transformer.TransformRequests(nil, records, testTransformation)

	var actual = chanToSet(requests)
	if len(actual) != 0 {
		t.Error("incorrect result: expected set size is 0, got", len(actual))
	}
}

func TestTransformRequests(t *testing.T) {
	records := make(chan reader.Record)
	go func() {
		records <- reader.Record{Values: []string{"/api/test1"}}
		records <- reader.Record{Values: []string{"/api/test2"}}
		close(records)
	}()

	requests := transformer.TransformRequests([]string{"http://test1.com", "http://test2.com"}, records, testTransformation)

	var actual = chanToSet(requests)
	if len(actual) != 4 {
		t.Error("incorrect result: expected set size is 4, got", len(actual))
	}

	expected := []sender.Request{
		{Url: "http://test1.com/api/test1", Method: "GET"},
		{Url: "http://test2.com/api/test1", Method: "GET"},
		{Url: "http://test1.com/api/test2", Method: "GET"},
		{Url: "http://test2.com/api/test2", Method: "GET"},
	}

	for _, req := range expected {
		_, ok := actual[req]
		if !ok {
			t.Errorf("incorrect result: expected request %v not found", req)
		}
	}
}

func TestTransformRequestsWithIncorrectRecords(t *testing.T) {
	records := make(chan reader.Record)
	go func() {
		records <- reader.Record{Values: []string{"/api/test1"}}
		records <- reader.Record{Values: []string{"/api/test2"}}
		close(records)
	}()

	requests := transformer.TransformRequests([]string{"http://test1.com", "http://test2.com"}, records, errorTransformation)

	var actual = chanToSet(requests)
	if len(actual) != 0 {
		t.Error("incorrect result: expected set size is 0, got", len(actual))
	}
}

func testTransformation(host string, rec reader.Record) (sender.Request, error) {
	return sender.Request{Url: host + rec.Values[0]}, nil
}

func errorTransformation(_ string, _ reader.Record) (sender.Request, error) {
	return sender.Request{}, errors.New("error")
}

func chanToSet(input <-chan sender.Request) map[sender.Request]bool {
	set := make(map[sender.Request]bool)
	for r := range input {
		set[r] = true
	}
	return set
}
