package reporter

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"testpoint/internal/comparator"
)

type CsvReporter struct {
	filename string
}

func NewCsvReporter(filename string) CsvReporter {
	return CsvReporter{filename: filename}
}

func (r CsvReporter) report(input <-chan comparator.RespDiff) {
	file := createFile(r.filename)
	defer file.Close()

	writer := csv.NewWriter(file)
	writeLine(writer, []string{
		"req_url_1", "req_url_2", "req_method", "req_headers", "req_body", "req_hash",
		"resp_status_1", "resp_body_1",
		"resp_status_2", "resp_body_2",
	})

	for d := range input {
		reqHash := strconv.FormatUint(d.Rec1.ReqHash, 10)
		writeLine(writer, []string{
			d.Rec1.ReqUrl, d.Rec2.ReqUrl, d.Rec1.ReqMethod, d.Rec1.ReqHeaders, d.Rec1.ReqBody, reqHash,
			d.Rec1.RespStatus, d.Rec1.RespBody,
			d.Rec2.RespStatus, d.Rec2.RespBody,
		})
	}

	log.Printf("the csv report was saved in %v", r.filename)
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
