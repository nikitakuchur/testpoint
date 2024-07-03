package writer

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testpoint/internal/sender"
)

// WriteResponses creates files for each unique host and writes the results in them.
func WriteResponses(input <-chan sender.RequestResponse, dir string) {
	fileMap := make(map[string]*os.File)
	writerMap := make(map[string]*csv.Writer)

	defer func() {
		for k, file := range fileMap {
			writerMap[k].Flush()
			closeFile(file)
		}
	}()

	for rr := range input {
		userUrl := rr.Request.UserUrl

		file, ok := fileMap[userUrl]
		writer := writerMap[userUrl]
		if !ok {
			path := filepath.Join(dir, urlToFilename(userUrl))
			file = createFile(path)

			fileMap[userUrl] = file
			writer = csv.NewWriter(file)
			writerMap[userUrl] = writer

			writeLine(writer, []string{
				"request_url", "request_method", "request_headers", "request_body",
				"response_status", "response_body",
			})
		}

		writeLine(writer, []string{
			rr.Request.Url, rr.Request.Method, rr.Request.Headers, rr.Request.Body,
			rr.Response.Status, rr.Response.Body,
		})
	}
}

func urlToFilename(url string) string {
	url = strings.ReplaceAll(url, "://", "-")
	url = strings.ReplaceAll(url, ":", "-")
	url = strings.ReplaceAll(url, "/", "-")
	url = strings.ReplaceAll(url, ".", "-")
	return url + ".csv"
}

func createFile(path string) *os.File {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalln("cannot create a new file:", err)
	}
	return file
}

func writeLine(writer *csv.Writer, record []string) {
	err := writer.Write(record)
	if err != nil {
		log.Fatalln("cannot write into a file:", err)
	}
}

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Fatalln("cannot close a file:", err)
	}
}
