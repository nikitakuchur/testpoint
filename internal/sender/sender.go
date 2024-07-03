package sender

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Request struct {
	Url     string
	Method  string
	Headers string
	Body    string
}

type Response struct {
	Status string
	Body   string
}

type RequestResponse struct {
	Request  Request
	Response Response
}

// SendRequests takes requests from the input channel, sends them to
// the corresponding endpoint, and puts the result in the output channel.
func SendRequests(input <-chan Request, workers int) <-chan RequestResponse {
	output := make(chan RequestResponse)

	client := &http.Client{}

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for req := range input {
				resp, err := sendRequest(client, req)
				if err != nil {
					log.Println(req, err)
					continue
				}
				output <- RequestResponse{req, resp}
			}
		}()
	}

	// this goroutine closes the channel
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

func sendRequest(client *http.Client, req Request) (Response, error) {
	httpReq, err := http.NewRequest(req.Method, req.Url, strings.NewReader(req.Body))
	if err != nil {
		return Response{}, errors.New("cannot create an http request: " + err.Error())
	}

	if req.Headers != "" {
		headersMap := map[string]string{}
		err = json.Unmarshal([]byte(req.Headers), &headersMap)
		if err != nil {
			return Response{}, errors.New("cannot convert headers to a map")
		}
		for k, v := range headersMap {
			httpReq.Header.Set(k, v)
		}
	}

	resp, err := doRequest(client, httpReq, 5)
	if err != nil {
		return Response{}, err
	}
	defer closeResponse(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, errors.New("cannot read an http body: " + err.Error())
	}

	return Response{resp.Status, string(body)}, nil
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
