package respwriter

import (
	"encoding/csv"
	"github.com/nikitakuchur/testpoint/internal/sender"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
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

	var processed atomic.Uint64
	ticker := time.NewTicker(10 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				log.Printf("collected %v responses...", processed.Load())
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
				"req_url", "req_method", "req_headers", "req_body", "req_hash",
				"resp_status", "resp_body",
			})
		}

		reqHash := strconv.FormatUint(rr.Request.Hash, 10)
		writeLine(writer, []string{
			rr.Request.Url, rr.Request.Method, rr.Request.Headers, rr.Request.Body, reqHash,
			rr.Response.Status, rr.Response.Body,
		})
		processed.Add(1)
	}
	ticker.Stop()
	done <- true
	log.Println("total number of collected responses:", processed.Load())
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
