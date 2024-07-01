package writer

import (
	"encoding/csv"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"testpoint/internal/sender"
)

// WriteResponses creates files for each unique host and writes the results in them.
func WriteResponses(input <-chan sender.Response, dir string) {
	fileMap := make(map[string]*os.File)
	writerMap := make(map[string]*csv.Writer)

	defer func() {
		for k, file := range fileMap {
			writerMap[k].Flush()
			closeFile(file)
		}
	}()

	for res := range input {
		u, err := url.Parse(res.Request.Url)
		if err != nil {
			log.Fatalln("cannot parse a URL:", err)
		}

		file, ok := fileMap[u.Host]
		writer := writerMap[u.Host]
		if !ok {
			path := filepath.Join(dir, u.Host+".csv")
			file = createFile(path)

			fileMap[u.Host] = file
			writer = csv.NewWriter(file)
			writerMap[u.Host] = writer

			writeLine(writer, []string{"url", "response"})
		}

		writeLine(writer, []string{res.Request.Url, res.Response})
	}
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
