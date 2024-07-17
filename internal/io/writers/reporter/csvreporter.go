package reporter

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"testpoint/internal/comparator"
)

// CsvReporter represents a reporter that writes mismatched records to a CSV file.
type CsvReporter struct {
	filename string
}

func NewCsvReporter(filename string) CsvReporter {
	return CsvReporter{filename: filename}
}

func (r CsvReporter) Report(input <-chan comparator.RespDiff) {
	file := createFile(r.filename)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writeLine(writer, []string{
		"req1_url", "req1_method", "req1_headers", "req1_body",
		"req2_url", "req2_method", "req2_headers", "req2_body",
		"req_hash",
		"resp1_status", "resp1_body",
		"resp2_status", "resp2_body",
	})

	for d := range input {
		reqHash := strconv.FormatUint(d.Rec1.ReqHash, 10)
		writeLine(writer, []string{
			d.Rec1.ReqUrl, d.Rec1.ReqMethod, d.Rec1.ReqHeaders, d.Rec1.ReqBody,
			d.Rec2.ReqUrl, d.Rec2.ReqMethod, d.Rec2.ReqHeaders, d.Rec2.ReqBody,
			reqHash,
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
