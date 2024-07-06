package respreader

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

type RespRecord struct {
	ReqUrl     string
	ReqMethod  string
	ReqHeaders string
	ReqBody    string
	ReqHash    uint64

	RespStatus string
	RespBody   string
}

// ReadResponses reads the CSV file with responses and sends the data to the output channel.
func ReadResponses(filename string) <-chan RespRecord {
	output := make(chan RespRecord)

	go func() {
		defer close(output)

		file, err := os.Open(filename)
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		err = readRecords(file, output)
		if err != nil {
			log.Fatalln("cannot process the given file", err)
		}
	}()

	return output
}

func readRecords(file *os.File, output chan<- RespRecord) error {
	reader := csv.NewReader(file)

	// we need to skip the CSV header
	header, err := reader.Read()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		log.Fatalf("%v: %v", file.Name(), err)
	}
	if len(header) < 7 {
		log.Fatalf("%v: there are missing values", file.Name())
	}

	for {
		values, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("%v: %v, the record was skipped", file.Name(), err)
			continue
		}

		hash, err := strconv.ParseUint(values[4], 10, 64)
		if err != nil {
			log.Printf("%v: cannot parse the hash value '%v', the record was skipped", file.Name(), values[4])
			continue
		}

		rec := RespRecord{
			values[0], values[1], values[2], values[3], hash,
			values[5], values[6],
		}
		output <- rec
	}
}