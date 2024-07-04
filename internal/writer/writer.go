package writer

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testpoint/internal/sender"
	"time"
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

	processed := 0
	ticker := time.NewTicker(10 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				log.Printf("processed %v requests...", processed)
			}
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
				"request_url", "request_method", "request_headers", "request_body", "request_hash",
				"response_status", "response_body",
			})
		}

		writeLine(writer, []string{
			rr.Request.Url, rr.Request.Method, rr.Request.Headers, rr.Request.Body, strconv.FormatUint(rr.Request.Hash, 10),
			rr.Response.Status, rr.Response.Body,
		})
		processed++
	}
	ticker.Stop()
	done <- true
	log.Println("total number of processed requests:", processed)
}

func urlToFilename(url string) string {
	if url == "" {
		return "output.csv"
	}
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
