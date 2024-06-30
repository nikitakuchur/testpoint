package reader

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Record struct {
	Header []string
	Values []string
}

// ReadRequests reads the CSV files and sends the data to the output channel.
func ReadRequests(path string, withHeader bool) <-chan Record {
	output := make(chan Record)

	go func() {
		defer close(output)

		if path == "" {
			return
		}

		filenames := readFilenames(path)

		for _, filename := range filenames {
			file, err := os.Open(filename)
			if err != nil {
				log.Fatalln("cannot read the given input file:", err)
			}

			reader := csv.NewReader(file)

			var header []string = nil
			if withHeader {
				header, err = reader.Read()
				if err != nil {
					log.Fatalln("cannot read the header of the input file")
				}
			}

			for {
				rec, err := reader.Read()
				if err == io.EOF {
					break
				}
				output <- Record{header, rec}
			}

			closeFile(file)
		}
	}()

	return output
}

func readFilenames(path string) []string {
	var filenames []string

	if !isDir(path) {
		filenames = append(filenames, path)
		return filenames
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatalln("cannot read the directory:", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filename := filepath.Join(path, entry.Name())
			filenames = append(filenames, filename)
		}
	}

	return filenames
}

func isDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		log.Fatalln("cannot read the given input files:", err)
	}
	return stat.IsDir()
}

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Fatalln("cannot close a file:", err)
	}
}
