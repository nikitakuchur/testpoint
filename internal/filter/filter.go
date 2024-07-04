package filter

import (
	"testpoint/internal/reader"
)

// Filter removes duplicates from the data stream
func Filter(input <-chan reader.Record) <-chan reader.Record {
	output := make(chan reader.Record)

	set := make(map[uint64]interface{})

	go func() {
		defer close(output)

		for rec := range input {
			_, ok := set[rec.Hash]
			if ok {
				continue
			}
			output <- rec
			set[rec.Hash] = struct{}{}
		}
	}()

	return output
}
