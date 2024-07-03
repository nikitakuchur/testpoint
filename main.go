package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"testpoint/internal/reader"
	"testpoint/internal/sender"
	"testpoint/internal/transformer"
	"testpoint/internal/writer"
)

type config struct {
	input     string
	header    bool
	urls      []string
	transform string
	workers   int
	output    string
}

func (c config) String() string {
	return fmt.Sprintf(
		"input: \"%v\", header: %v, urls: %v, transform: %v, workers: %v, output: \"%v\"",
		c.input, c.header, c.urls, c.transform, c.workers, c.output,
	)
}

func main() {
	inputPtr := flag.String("input", "", "a CSV file or directory with CSV files")
	headerPtr := flag.Bool("no-header", true, "enable this flag if your CSV file has no header")
	hostsPtr := flag.String("urls", "", "a list of hosts to send requests to")
	transformPtr := flag.String("transform", "", "a JavaScript file with a request transformation")
	workPtr := flag.Int("w", 1, "a number of workers to send requests")
	outputPtr := flag.String("output", "./", "a directory where the output files need to be saved")

	flag.Parse()

	urls := parseUrls(*hostsPtr)

	conf := config{
		*inputPtr,
		*headerPtr,
		urls,
		*transformPtr,
		*workPtr,
		*outputPtr,
	}

	log.Println(conf)
	log.Println("starting to process the requests...")

	records := reader.ReadRequests(conf.input, conf.header)

	requests := transformer.TransformRequests(urls, records, createTransformation(conf.transform))
	responses := sender.SendRequests(requests, conf.workers)
	writer.WriteResponses(responses, conf.output)

	log.Println("completed")
	log.Printf("the result is saved in %v", conf.output)
}

func createTransformation(filepath string) transformer.Transformation {
	if filepath == "" {
		return transformer.DefaultTransformation
	}
	script := readScript(filepath)
	transformation, err := transformer.NewTransformation(script)
	if err != nil {
		log.Fatalln(err)
	}
	return transformation
}

func parseUrls(urls string) []string {
	return strings.Split(strings.ReplaceAll(urls, " ", ""), ",")
}

func readScript(filename string) string {
	script, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln("cannot read the transformation script:", err)
	}
	return string(script)
}
