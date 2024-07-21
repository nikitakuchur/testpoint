package sender_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/nikitakuchur/testpoint/internal/sender"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendRequestsWithNoRequests(t *testing.T) {
	requests := make(chan sender.Request)
	close(requests)

	responses := sender.SendRequests(requests, 1)

	actual := chanToSlice(responses)

	if len(actual) != 0 {
		t.Error("incorrect result: expected number of responses is 0, got", len(actual))
	}
}

func TestSendRequestsWithZeroWorkers(t *testing.T) {
	requests := make(chan sender.Request)
	go func() {
		requests <- sender.Request{
			Url:    "http://test.com/api/test",
			Method: "GET",
		}
		close(requests)
	}()

	responses := sender.SendRequests(requests, 0)

	actual := chanToSlice(responses)

	if len(actual) != 0 {
		t.Error("incorrect result: expected number of responses is 0, got", len(actual))
	}
}

func TestSendRequests(t *testing.T) {
	handlerFunc := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, err := rw.Write([]byte("Hello world!"))
		if err != nil {
			log.Fatalln("cannot write a response body")
		}
	})
	server := httptest.NewServer(handlerFunc)
	defer server.Close()

	requests := make(chan sender.Request)
	go func() {
		requests <- sender.Request{
			Url:     server.URL,
			Method:  "GET",
			Headers: `{"myHeader":"foo"}`,
		}
		close(requests)
	}()

	responses := sender.SendRequests(requests, 1)

	actual := chanToSlice(responses)

	if len(actual) != 1 {
		t.Error("incorrect result: expected number of responses is 1, got", len(actual))
	}
	expected := []sender.RequestResponse{
		{
			sender.Request{Url: server.URL, Method: "GET", Headers: `{"myHeader":"foo"}`},
			sender.Response{Status: "200", Body: "Hello world!"},
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Error(diff)
	}
}

func TestSendRequestsWithIncorrectRequest(t *testing.T) {
	requests := make(chan sender.Request)
	go func() {
		requests <- sender.Request{
			Url:    ":http://test.com/api/test",
			Method: "GET",
		}
		close(requests)
	}()

	responses := sender.SendRequests(requests, 1)

	actual := chanToSlice(responses)

	if len(actual) != 0 {
		t.Error("incorrect result: expected number of responses is 0, got", len(actual))
	}
}

func TestSendRequestsWithIncorrectHeaders(t *testing.T) {
	requests := make(chan sender.Request)
	go func() {
		requests <- sender.Request{
			Url:     "http://test.com/api/test",
			Method:  "GET",
			Headers: "123",
		}
		close(requests)
	}()

	responses := sender.SendRequests(requests, 1)

	actual := chanToSlice(responses)

	if len(actual) != 0 {
		t.Error("incorrect result: expected number of responses is 0, got", len(actual))
	}
}

func chanToSlice(input <-chan sender.RequestResponse) []sender.RequestResponse {
	var slice []sender.RequestResponse
	for rec := range input {
		slice = append(slice, rec)
	}
	return slice
}
