package reader

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ReadRequests reads the CSV files and sends the data into the given channel.
func ReadRequests(path string, output chan<- []string) {
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

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			row := strings.Split(scanner.Text(), ",")
			output <- row
		}

		closeFile(file)
	}
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
