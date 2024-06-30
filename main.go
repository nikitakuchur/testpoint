package main

import (
	"flag"
	"fmt"
	"log"
	"restcompare/internal/reader"
	"restcompare/internal/sender"
	"restcompare/internal/transformer"
	"restcompare/internal/writer"
	"strings"
)

type config struct {
	input   string
	header  bool
	workers int
	hosts   []string
	output  string
}

func (c config) String() string {
	return fmt.Sprintf(
		"intput: \"%v\", header: %v, workers: %v, hosts: %v, output: \"%v\"",
		c.input, c.header, c.workers, c.hosts, c.output,
	)
}

func main() {
	inputPtr := flag.String("input", "", "a CSV file or directory with CSV files")
	headerPtr := flag.Bool("header", false, "enable this flag is your CSV file has a header")
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

	rowCh := reader.ReadRequests(conf.input)
	requestCh := transformer.TransformRequests(hosts, transform, rowCh)
	responseCh := sender.SendRequests(requestCh)
	writer.WriteResponses(responseCh, conf.output)

	log.Println("completed")
	log.Printf("the result is saved in %v", conf.output)
}

func transform(url string, row []string) transformer.Request {
	return transformer.Request{Url: url + row[1], Method: row[0]}
}

func parseHosts(hosts string) []string {
	return strings.Split(strings.ReplaceAll(hosts, " ", ""), ",")
}
