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
	workers    int
	hosts      []string
	output     string
}

func (c config) String() string {
	return fmt.Sprintf(
		"intput: \"%v\", header: %v, workers: %v, hosts: %v, output: \"%v\"",
		c.input, c.withHeader, c.workers, c.hosts, c.output,
	)
}

func main() {
	inputPtr := flag.String("input", "", "a CSV file or directory with CSV files")
	headerPtr := flag.Bool("header", false, "enable this flag if your CSV file has a header")
	workPtr := flag.Int("w", 1, "a number of workers to send requests")
	hostsPtr := flag.String("hosts", "", "a list of hosts to send requests to")
	outputPtr := flag.String("output", "./", "a directory where the output files need to be saved")

	flag.Parse()

	hosts := parseHosts(*hostsPtr)

	conf := config{
		*inputPtr,
		*headerPtr,
		*workPtr,
		hosts,
		*outputPtr,
	}

	log.Println(conf)
	log.Println("starting to process the requests...")

	records := reader.ReadRequests(conf.input, conf.withHeader)
	requests := transformer.TransformRequests(hosts, records, transformer.DefaultTransformation)
	responses := sender.SendRequests(requests, conf.workers)
	writer.WriteResponses(responses, conf.output)

	log.Println("completed")
	log.Printf("the result is saved in %v", conf.output)
}

func parseHosts(hosts string) []string {
	return strings.Split(strings.ReplaceAll(hosts, " ", ""), ",")
}
