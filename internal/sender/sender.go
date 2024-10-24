package sender

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Request struct {
	Url     string
	Method  string
	Headers string
	Body    string

	UserUrl string
	Hash    uint64
}

func (r Request) String() string {
	return fmt.Sprintf("{%v %v %v %v}", r.Url, r.Method, r.Headers, r.Body)
}

type Response struct {
	Status string
	Body   string
}

type RequestResponse struct {
	Request  Request
	Response Response
}

type Sender struct {
	client *http.Client
}

func NewSender() Sender {
	return Sender{&http.Client{}}
}

// SendRequests takes requests from the input channel, sends them to
// the corresponding endpoint, and puts the result in the output channel.
func (s Sender) SendRequests(input <-chan Request, workers int) <-chan RequestResponse {
	output := make(chan RequestResponse)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for req := range input {
				resp, err := s.sendRequest(req)
				if err != nil {
					log.Printf("%v: %v, request was skipped", req, err)
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

func (s Sender) sendRequest(req Request) (Response, error) {
	httpReq, err := http.NewRequest(req.Method, req.Url, strings.NewReader(req.Body))
	if err != nil {
		return Response{}, fmt.Errorf("cannot create an http request: %w", err)
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

	resp, err := s.doRequest(httpReq, 5)
	if err != nil {
		return Response{}, err
	}
	defer closeResponse(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("cannot read an http body: %w", err)
	}

	status := strconv.FormatInt(int64(resp.StatusCode), 10)
	return Response{status, string(body)}, nil
}

func (s Sender) doRequest(req *http.Request, retries int) (*http.Response, error) {
	for i := 0; i < retries; i++ {
		resp, err := s.client.Do(req)
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
