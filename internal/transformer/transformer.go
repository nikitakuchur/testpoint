package transformer

import (
	"log"
	"testpoint/internal/reader"
	"testpoint/internal/sender"
)

type Transformation func(host string, rec reader.Record) (sender.Request, error)

// TransformRequests reads raw request data from the input channel,
// transforms it into requests using the given transformation and sends it to the output channel.
func TransformRequests(hosts []string, input <-chan reader.Record, transformation Transformation) <-chan sender.Request {
	output := make(chan sender.Request)

	go func() {
		defer close(output)

		if len(hosts) == 0 {
			return
		}

	outer:
		for rec := range input {
			for _, host := range hosts {
				req, err := transformation(host, rec)
				if err != nil {
					log.Printf("%v: %v, the record was skipped", rec, err)
					continue outer
				}
				if req.Method == "" {
					req.Method = "GET"
				}
				output <- req
			}
		}
	}()

	return output
}
