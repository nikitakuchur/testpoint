package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"testpoint/internal/reader"
	"testpoint/internal/sender"
	"testpoint/internal/transformer"
	"testpoint/internal/writer"
)

type config struct {
	input      string
	withHeader bool
	hosts      []string
	transform  string
	workers    int
	output     string
}

func (c config) String() string {
	return fmt.Sprintf(
		"intput: \"%v\", header: %v, hosts: %v, transform: %v, workers: %v, output: \"%v\"",
		c.input, c.withHeader, c.hosts, c.transform, c.workers, c.output,
	)
}

func main() {
	inputPtr := flag.String("input", "", "a CSV file or directory with CSV files")
	headerPtr := flag.Bool("header", false, "enable this flag if your CSV file has a header")
	hostsPtr := flag.String("hosts", "", "a list of hosts to send requests to")
	transformPtr := flag.String("transform", "", "a JavaScript file with a request transformation")
	workPtr := flag.Int("w", 1, "a number of workers to send requests")
	outputPtr := flag.String("output", "./", "a directory where the output files need to be saved")

	flag.Parse()

	hosts := parseHosts(*hostsPtr)

	conf := config{
		*inputPtr,
		*headerPtr,
		hosts,
		*transformPtr,
		*workPtr,
		*outputPtr,
	}

	log.Println(conf)
	log.Println("starting to process the requests...")

	records := reader.ReadRequests(conf.input, conf.withHeader)

	transformation := transformer.DefaultTransformation
	if conf.transform != "" {
		transformation = transformer.NewTransformation(conf.transform)
	}

	requests := transformer.TransformRequests(hosts, records, transformation)
	responses := sender.SendRequests(requests, conf.workers)
	writer.WriteResponses(responses, conf.output)

	log.Println("completed")
	log.Printf("the result is saved in %v", conf.output)
}

func parseHosts(hosts string) []string {
	return strings.Split(strings.ReplaceAll(hosts, " ", ""), ",")
}
