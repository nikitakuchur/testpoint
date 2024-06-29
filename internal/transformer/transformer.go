package transformer

type Request struct {
	Url     string
	Method  string
	Headers string
	Body    string
}

// TransformRequests reads from the input channel raw request data,
// transforms it into requests using the given transformer and sends it to the output channel
func TransformRequests(urls []string, transformer func(string, []string) Request, input <-chan []string, output chan<- Request) {
	for row := range input {
		for _, url := range urls {
			req := transformer(url, row)
			output <- req
		}
	}
	// TODO: we can't have more than one goroutines to run this function
	// 	use waitGroups to fix this issue
	close(output)
}
