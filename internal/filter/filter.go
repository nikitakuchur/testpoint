package filter

import (
	"testpoint/internal/reqreader"
)

// Filter removes duplicates from the data stream
func Filter(input <-chan reqreader.Record) <-chan reqreader.Record {
	output := make(chan reqreader.Record)

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
