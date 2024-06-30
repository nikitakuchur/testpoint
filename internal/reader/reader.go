package reader

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
)

// ReadRequests reads the CSV files and sends the data to the output channel.
func ReadRequests(path string) <-chan []string {
	output := make(chan []string)

	go func() {
		defer close(output)

		if path == "" {
			return
		}

		var filenames []string
		if isDir(path) {
			filenames = readFilenames(path)
		} else {
			filenames = append(filenames, path)
		}

		for _, filename := range filenames {
			file, err := os.Open(filename)
			if err != nil {
				log.Fatalln("cannot read the given input file:", err)
			}

			reader := csv.NewReader(file)
			for {
				rec, err := reader.Read()
				if err == io.EOF {
					break
				}
				output <- rec
			}

			closeFile(file)
		}
	}()

	return output
}

func readFilenames(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalln("cannot read the directory:", err)
	}
	var filenames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			filename := filepath.Join(dir, entry.Name())
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
