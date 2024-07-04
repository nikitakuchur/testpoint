package filter

import (
	"hash/fnv"
	"testpoint/internal/reader"
)

// Filter removes duplicates from the data stream
func Filter(input <-chan reader.Record) <-chan reader.Record {
	output := make(chan reader.Record)

	set := make(map[uint64]interface{})

	go func() {
		defer close(output)

		for rec := range input {
			h := hash(rec)
			_, ok := set[h]
			if ok {
				continue
			}
			output <- rec
			set[h] = struct{}{}
		}
	}()

	return output
}

func hash(rec reader.Record) uint64 {
	h := fnv.New64()
	h.Write([]byte(rec.String()))
	return h.Sum64()
}
