package main

import (
	"flag"
	"fmt"
	"log"
	"restcompare/internal/reader"
	"restcompare/internal/sender"
	"restcompare/internal/transformer"
	"strings"
)

type config struct {
	input   string
	header  bool
	workers int
	urls    []string
	output  string
}

func (c config) String() string {
	return fmt.Sprintf(
		"intput: %v, header: %v, workers: %v, urls: %v, output: %v",
		c.input, c.header, c.workers, c.urls, c.output,
	)
}

func main() {
	inputPtr := flag.String("input", "", "a CSV file or directory with CSV files")
	headerPtr := flag.Bool("header", false, "enable this flag is your CSV file has a header")
	workPtr := flag.Int("w", 1, "a number of workers to send requests")
	urlsPtr := flag.String("urls", "", "a list of URLs to send requests.")
	outputPtr := flag.String("output", "/output", "a directory where the output files need to be saved")

	flag.Parse()

	urls := parseUrls(*urlsPtr)

	conf := config{
		*inputPtr,
		*headerPtr,
		*workPtr,
		urls,
		*outputPtr,
	}

	log.Println(conf)
	log.Println("starting to process the requests...")

	rowCh := make(chan []string)
	go reader.ReadRequests(*inputPtr, rowCh)

	requestCh := make(chan transformer.Request)
	go transformer.TransformRequests(urls, transform, rowCh, requestCh)

	responseCh := make(chan sender.Response)
	go sender.SendRequests(requestCh, responseCh)

	// write everything into a file
	for res := range responseCh {
		log.Println(res)
	}
}

func transform(url string, row []string) transformer.Request {
	return transformer.Request{Url: url + row[1], Method: row[0]}
}

func parseUrls(urls string) []string {
	return strings.Split(strings.ReplaceAll(urls, " ", ""), ",")
}
