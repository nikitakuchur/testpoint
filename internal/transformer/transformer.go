package transformer

import "testpoint/internal/reader"

type Request struct {
	Url     string
	Method  string
	Headers map[string]string
	Body    string
}

type Transformation func(string, reader.Record) Request

// TransformRequests reads raw request data from the input channel,
// transforms it into requests using the given transformation and sends it to the output channel.
func TransformRequests(hosts []string, transformation Transformation, input <-chan reader.Record) <-chan Request {
	output := make(chan Request)

	go func() {
		defer close(output)

		if len(hosts) == 0 {
			return
		}

		for rec := range input {
			for _, url := range hosts {
				req := transformation(url, rec)
				output <- req
			}
		}
	}()

	return output
}
