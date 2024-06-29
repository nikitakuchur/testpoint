package sender

import "restcompare/internal/transformer"

type Response struct {
	Request  transformer.Request
	Response string
}

func SendRequests(input <-chan transformer.Request, output chan<- Response) {
	for req := range input {
		// send an http request
		output <- Response{req, "test response"}
	}
	close(output)
}
