package transformer

type Request struct {
	Url     string
	Method  string
	Headers string
	Body    string
}

// TransformRequests reads from the input channel raw request data,
// transforms it into requests using the given transformer and sends it to the output channel
func TransformRequests(hosts []string, transformer func(string, []string) Request, input <-chan []string, output chan<- Request) {
	// TODO: we can't have more than one goroutines to run this function
	// 	use waitGroups to fix this issue
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
}
