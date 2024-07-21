package filter

import (
	"github.com/nikitakuchur/testpoint/internal/io/readers/reqreader"
)

// Filter removes duplicates from the data stream.
func Filter(input <-chan reqreader.ReqRecord) <-chan reqreader.ReqRecord {
	output := make(chan reqreader.ReqRecord)

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
