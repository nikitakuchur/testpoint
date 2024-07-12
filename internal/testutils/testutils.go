package testutils

import (
	"log"
	"os"
)

// ChanToSlice collects all the data from the given channel and converts it into a slice.
func ChanToSlice[T any](input <-chan T) []T {
	var slice []T
	for rec := range input {
		slice = append(slice, rec)
	}
	return slice
}

// CreateTempFile creates a temp file with the given name and content.
func CreateTempFile(dir string, name string, content string) string {
	file, err := os.CreateTemp(dir, name)
	if err != nil {
		log.Fatalln("cannot create a temp file")
	}
	_, err = file.WriteString(content)
	if err != nil {
		log.Fatalln("cannot write into a temp file")
	}
	err = file.Close()
	if err != nil {
		log.Fatalln("cannot close a temp file")
	}
	return file.Name()
}

// ReadFile reads a the given file.
func ReadFile(filepath string) string {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("cannot read a file")
	}
	return string(bytes)
}
