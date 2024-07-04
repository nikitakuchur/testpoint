package transformer

import (
	"log"
	"testpoint/internal/reader"
	"testpoint/internal/sender"
)

type Transformation func(userUrls string, rec reader.Record) (sender.Request, error)

// TransformRequests reads raw request data from the input channel,
// transforms it into requests using the given transformation and sends it to the output channel.
func TransformRequests(userUrls []string, input <-chan reader.Record, transformation Transformation) <-chan sender.Request {
	output := make(chan sender.Request)

	go func() {
		defer close(output)

		if len(userUrls) == 0 {
			return
		}

		for rec := range input {
			for _, url := range userUrls {
				req, err := transformation(url, rec)
				if err != nil {
					log.Printf("%v, %v: %v, the record was skipped", url, rec, err)
					continue
				}

				if req.Method == "" {
					req.Method = "GET"
				}
				req.UserUrl = url
				req.Hash = rec.Hash

				output <- req
			}
		}
	}()

	return output
}
