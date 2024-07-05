package respreader

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

type Record struct {
	ReqUrl     string
	ReqMethod  string
	ReqHeaders string
	ReqBody    string
	ReqHash    uint64

	RespStatus string
	RespBody   string
}

// ReadResponses reads the CSV file with responses and sends the data to the output channel.
func ReadResponses(filename string) <-chan Record {
	output := make(chan Record)

	go func() {
		defer close(output)

		if filename == "" {
			return
		}

		file, err := os.Open(filename)
		if err != nil {
			log.Fatalln("cannot read the given file", err)
		}
		defer closeFile(file)

		err = readRecords(file, output)
		if err != nil {
			log.Fatalln("cannot process the given file", err)
		}
	}()

	return output
}

func readRecords(file *os.File, output chan<- Record) error {
	reader := csv.NewReader(file)

	// we need to skip the CSV header
	_, err := reader.Read()
	if err != nil {
		return err
	}

	for {
		values, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		hash, err := strconv.ParseUint(values[4], 10, 64)
		if err != nil {
			return err
		}

		rec := Record{
			values[0], values[1], values[2], values[3], hash,
			values[5], values[6],
		}
		output <- rec
	}
}

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Fatalln("cannot close a file:", err)
	}
}
