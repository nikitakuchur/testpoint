package reqreader

import (
	"encoding/csv"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ReqRecord struct {
	Fields []string
	Values []string
	Hash   uint64
}

func (rec ReqRecord) String() string {
	if rec.Fields != nil {
		var sb strings.Builder
		for i, field := range rec.Fields {
			sb.WriteString(fmt.Sprintf("%v: %v", field, rec.Values[i]))
			if i != len(rec.Fields)-1 {
				sb.WriteString(", ")
			}
		}
		return sb.String()
	}
	return strings.Join(rec.Values, ", ")
}

// ReadRequests reads the CSV files with requests and sends the data to the output channel.
func ReadRequests(path string, withHeader bool) <-chan ReqRecord {
	output := make(chan ReqRecord)

	go func() {
		defer close(output)

		filenames, err := readFilenames(path)
		if err != nil {
			log.Printf("%v: %v, request reading was skipped", path, err)
			return
		}

		for _, filename := range filenames {
			err := readFile(filename, withHeader, output)
			if err != nil {
				log.Printf("%v: %v, the file was skipped", filename, err)
			}
		}
	}()

	return output
}

func readFile(filename string, withHeader bool, output chan<- ReqRecord) error {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return err
	}

	err = readRecords(file, withHeader, output)
	if err != nil {
		return err
	}

	return nil
}

func readRecords(file *os.File, withHeader bool, output chan<- ReqRecord) error {
	reader := csv.NewReader(file)

	var header []string = nil
	if withHeader {
		h, err := reader.Read()
		if err != nil {
			return err
		}
		header = h
	}

	for {
		values, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("%v, the record was skipped", err)
			continue
		}

		rec := ReqRecord{Fields: header, Values: values}
		rec.Hash = hash(rec)
		output <- rec
	}
}

func hash(rec ReqRecord) uint64 {
	h := fnv.New64()
	h.Write([]byte(rec.String()))
	return h.Sum64()
}

func readFilenames(path string) ([]string, error) {
	dir, err := isDir(path)
	if err != nil {
		return nil, err
	}

	if !dir {
		return []string{path}, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			filename := filepath.Join(path, entry.Name())
			filenames = append(filenames, filename)
		}
	}

	return filenames, nil
}

func isDir(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return stat.IsDir(), nil
}
