package sender

import (
	"errors"
	"io"
	"log"
	"net/http"
	"restcompare/internal/transformer"
	"strings"
	"time"
)

type Response struct {
	Request  transformer.Request
	Response string
}

// SendRequests takes requests from the input channel, sends them to
// the corresponding host, and puts the result in the output channel.
func SendRequests(input <-chan transformer.Request) <-chan Response {
	output := make(chan Response)

	go func() {
		client := &http.Client{}

		for req := range input {
			body, err := sendRequest(client, req)
			if err != nil {
				log.Println(req, err)
				continue
			}
			output <- Response{req, body}
		}

		close(output)
	}()

	return output
}

func sendRequest(client *http.Client, req transformer.Request) (string, error) {
	httpReq, err := http.NewRequest(req.Method, req.Url, strings.NewReader(req.Body))
	if err != nil {
		return "", errors.New("cannot create an http request: " + err.Error())
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	resp, err := doRequest(client, httpReq, 5)
	if err != nil {
		return "", err
	}
	defer closeResponse(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("cannot read an http body: " + err.Error())
	}

	return string(body), nil
}

func doRequest(client *http.Client, req *http.Request, retries int) (*http.Response, error) {
	for i := 0; i < retries; i++ {
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("%v, retry attempt=%v", err, i+1)
			time.Sleep(2 * time.Second)
			continue
		}
		return resp, nil
	}
	return nil, errors.New("cannot send an http request")
}

func closeResponse(resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		log.Fatalln("cannot close a response body:", err)
	}
}
