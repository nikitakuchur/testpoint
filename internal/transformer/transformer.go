package transformer

import (
	"log"
	"testpoint/internal/io/readers/reqreader"
	"testpoint/internal/sender"
)

// TransformRequests reads raw request data from the input channel,
// transforms it into requests using the given transformation and sends it to the output channel.
func TransformRequests(userUrls []string, input <-chan reqreader.ReqRecord, transformation Transformation) <-chan sender.Request {
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
