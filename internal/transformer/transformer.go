package transformer

type Request struct {
	Url     string
	Method  string
	Headers string
	Body    string
}

// TransformRequests reads from the input channel raw request data,
// transforms it into requests using the given transformer and sends it to the output channel.
func TransformRequests(hosts []string, transformer func(string, []string) Request, input <-chan []string) <-chan Request {
	output := make(chan Request)

	go func() {
		defer close(output)

		if len(hosts) == 0 {
			return
		}

		for row := range input {
			for _, url := range hosts {
				req := transformer(url, row)
				output <- req
			}
		}
	}()

	return output
}
