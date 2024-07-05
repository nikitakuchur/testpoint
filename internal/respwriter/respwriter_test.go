package respwriter_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"testpoint/internal/respwriter"
	"testpoint/internal/sender"
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
				Headers: `{"my_header":"foo"}`,
				Body:    `{"field":"foo"}`,
				UserUrl: "http://test.com",
			},
			Response: sender.Response{Status: "200", Body: "Hello world!"},
		}
		responses <- sender.RequestResponse{
			Request: sender.Request{
				Url:     "http://test.com/api/bar",
				Method:  "GET",
				Headers: `{"my_header":"bar"}`,
				Body:    `{"field":"bar"}`,
				UserUrl: "http://test.com",
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

	actual := readFile(tempDir + "/http-test-com.csv")

	expected := `request_url,request_method,request_headers,request_body,response_status,response_body
http://test.com/api/foo,GET,"{""my_header"":""foo""}","{""field"":""foo""}",200,Hello world!
http://test.com/api/bar,GET,"{""my_header"":""bar""}","{""field"":""bar""}",200,Goodbye!
`

	if actual != expected {
		t.Errorf("incorrect result:\n expected: %v\nactual: %v", expected, actual)
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
				Headers: `{"my_header":"foo"}`,
				Body:    `{"field":"foo"}`,
				UserUrl: "http://test1.com",
			},
			Response: sender.Response{Status: "200", Body: "Hello world!"},
		}
		responses <- sender.RequestResponse{
			Request: sender.Request{
				Url:     "http://test2.com/api/bar",
				Method:  "GET",
				Headers: `{"my_header":"bar"}`,
				Body:    `{"field":"bar"}`,
				UserUrl: "http://test2.com",
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
		{"/http-test1-com.csv", `request_url,request_method,request_headers,request_body,response_status,response_body
http://test1.com/api/foo,GET,"{""my_header"":""foo""}","{""field"":""foo""}",200,Hello world!
`},
		{"/http-test2-com.csv", `request_url,request_method,request_headers,request_body,response_status,response_body
http://test2.com/api/bar,GET,"{""my_header"":""bar""}","{""field"":""bar""}",200,Goodbye!
`},
	}

	for _, e := range expected {
		actual := readFile(tempDir + e.filename)

		if actual != e.content {
			t.Errorf("incorrect result:\n expected: %v\nactual: %v", e.content, actual)
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
				Headers: `{"my_header":"foo"}`,
				Body:    `{"field":"foo"}`,
				UserUrl: "",
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

	actual := readFile(tempDir + "/output.csv")

	expected := `request_url,request_method,request_headers,request_body,response_status,response_body
http://test.com/api/foo,GET,"{""my_header"":""foo""}","{""field"":""foo""}",200,Hello world!
`

	if actual != expected {
		t.Errorf("incorrect result:\n expected: %v\nactual: %v", expected, actual)
	}
}

func readFile(filepath string) string {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("cannot read a file")
	}
	return string(bytes)
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
