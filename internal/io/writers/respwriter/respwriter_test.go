package respwriter_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"testpoint/internal/io/writers/respwriter"
	"testpoint/internal/sender"
	testutils "testpoint/internal/utils/testing"
)

func TestWriteResponsesWithNoResponses(t *testing.T) {
	tempDir := t.TempDir()

	responses := make(chan sender.RequestResponse)
	close(responses)

	respwriter.WriteResponses(responses, tempDir)

	filenames := readFilenames(tempDir)

	if len(filenames) != 0 {
		t.Error("incorrect result: expected number of files is 0, got", len(filenames))
	}
}

func TestWriteResponsesToOneFile(t *testing.T) {
	tempDir := t.TempDir()

	responses := make(chan sender.RequestResponse)
	go func() {
		responses <- sender.RequestResponse{
			Request: sender.Request{
				Url:     "http://test.com/api/foo",
				Method:  "GET",
				Headers: `{"myHeader":"foo"}`,
				Body:    `{"field":"foo"}`,
				UserUrl: "http://test.com",
				Hash:    1234,
			},
			Response: sender.Response{Status: "200", Body: "Hello world!"},
		}
		responses <- sender.RequestResponse{
			Request: sender.Request{
				Url:     "http://test.com/api/bar",
				Method:  "GET",
				Headers: `{"myHeader":"bar"}`,
				Body:    `{"field":"bar"}`,
				UserUrl: "http://test.com",
				Hash:    5678,
			},
			Response: sender.Response{Status: "200", Body: "Goodbye!"},
		}
		close(responses)
	}()

	respwriter.WriteResponses(responses, tempDir)

	filenames := readFilenames(tempDir)

	if len(filenames) != 1 {
		t.Error("incorrect result: expected number of files is 1, got", len(filenames))
	}

	actual := testutils.ReadFile(tempDir + "/http-test-com.csv")

	expected := `req_url,req_method,req_headers,req_body,req_hash,resp_status,resp_body
http://test.com/api/foo,GET,"{""myHeader"":""foo""}","{""field"":""foo""}",1234,200,Hello world!
http://test.com/api/bar,GET,"{""myHeader"":""bar""}","{""field"":""bar""}",5678,200,Goodbye!
`

	if actual != expected {
		t.Errorf("incorrect result:\nexpected: %v\nactual: %v", expected, actual)
	}
}

func TestWriteResponsesToMultipleFiles(t *testing.T) {
	tempDir := t.TempDir()

	responses := make(chan sender.RequestResponse)
	go func() {
		responses <- sender.RequestResponse{
			Request: sender.Request{
				Url:     "http://test1.com/api/foo",
				Method:  "GET",
				Headers: `{"myHeader":"foo"}`,
				Body:    `{"field":"foo"}`,
				UserUrl: "http://test1.com",
				Hash:    1234,
			},
			Response: sender.Response{Status: "200", Body: "Hello world!"},
		}
		responses <- sender.RequestResponse{
			Request: sender.Request{
				Url:     "http://test2.com/api/bar",
				Method:  "GET",
				Headers: `{"myHeader":"bar"}`,
				Body:    `{"field":"bar"}`,
				UserUrl: "http://test2.com",
				Hash:    5678,
			},
			Response: sender.Response{Status: "200", Body: "Goodbye!"},
		}
		close(responses)
	}()

	respwriter.WriteResponses(responses, tempDir)

	filenames := readFilenames(tempDir)

	if len(filenames) != 2 {
		t.Error("incorrect result: expected number of files is 2, got", len(filenames))
	}

	expected := []struct {
		filename string
		content  string
	}{
		{"/http-test1-com.csv", `req_url,req_method,req_headers,req_body,req_hash,resp_status,resp_body
http://test1.com/api/foo,GET,"{""myHeader"":""foo""}","{""field"":""foo""}",1234,200,Hello world!
`},
		{"/http-test2-com.csv", `req_url,req_method,req_headers,req_body,req_hash,resp_status,resp_body
http://test2.com/api/bar,GET,"{""myHeader"":""bar""}","{""field"":""bar""}",5678,200,Goodbye!
`},
	}

	for _, e := range expected {
		actual := testutils.ReadFile(tempDir + e.filename)

		if actual != e.content {
			t.Errorf("incorrect result:\nexpected: %v\nactual: %v", e.content, actual)
		}
	}
}

func TestWriteResponsesWithNoUrl(t *testing.T) {
	tempDir := t.TempDir()

	responses := make(chan sender.RequestResponse)
	go func() {
		responses <- sender.RequestResponse{
			Request: sender.Request{
				Url:     "http://test.com/api/foo",
				Method:  "GET",
				Headers: `{"myHeader":"foo"}`,
				Body:    `{"field":"foo"}`,
				UserUrl: "",
				Hash:    1234,
			},
			Response: sender.Response{Status: "200", Body: "Hello world!"},
		}
		close(responses)
	}()

	respwriter.WriteResponses(responses, tempDir)

	filenames := readFilenames(tempDir)

	if len(filenames) != 1 {
		t.Error("incorrect result: expected number of files is 1, got", len(filenames))
	}

	actual := testutils.ReadFile(tempDir + "/output.csv")

	expected := `req_url,req_method,req_headers,req_body,req_hash,resp_status,resp_body
http://test.com/api/foo,GET,"{""myHeader"":""foo""}","{""field"":""foo""}",1234,200,Hello world!
`

	if actual != expected {
		t.Errorf("incorrect result:\nexpected: %v\nactual: %v", expected, actual)
	}
}

func readFilenames(path string) []string {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("cannot read a directory")
	}

	var filenames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			filename := filepath.Join(path, entry.Name())
			filenames = append(filenames, filename)
		}
	}

	return filenames
}
