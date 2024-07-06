package send

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"testpoint/internal/filter"
	"testpoint/internal/io/readers/reqreader"
	"testpoint/internal/io/writers/respwriter"
	"testpoint/internal/sender"
	"testpoint/internal/transformer"
)

type config struct {
	input     string
	noHeader  bool
	urls      []string
	transform string
	workers   int
	output    string
}

func (c config) String() string {
	transform := c.transform
	if transform == "" {
		transform = "default"
	}
	return fmt.Sprintf(
		"input: '%v', noHeader: %v, urls: %v, transform: %v, workers: %v, output: '%v'",
		c.input, c.noHeader, c.urls, transform, c.workers, c.output,
	)
}

func Command() {
	inputPtr := flag.String("input", "", "a CSV file or directory with CSV files")
	noHeaderPtr := flag.Bool("no-header", false, "enable this flag if your CSV file has no header")
	hostsPtr := flag.String("urls", "", "a list of hosts, separated by commas, to which requests are to be sent")
	transformPtr := flag.String("transform", "", "a JavaScript file with a request transformation")
	workPtr := flag.Int("w", 1, "a number of workers to send requests")
	outputPtr := flag.String("output", "./", "a directory where the output files need to be saved")

	flag.Parse()

	conf := config{
		*inputPtr,
		*noHeaderPtr,
		parseUrls(*hostsPtr),
		*transformPtr,
		*workPtr,
		*outputPtr,
	}

	log.Printf("configuration: {%v}\n", conf)
	log.Println("starting to process the requests...")

	// TODO: replace it with a mandatory argument
	if conf.input == "" {
		log.Fatalln("input has to be specified")
	}

	records := reqreader.ReadRequests(conf.input, !conf.noHeader)

	records = filter.Filter(records)

	requests := transformer.TransformRequests(conf.urls, records, createTransformation(conf.transform))
	responses := sender.SendRequests(requests, conf.workers)
	respwriter.WriteResponses(responses, conf.output)

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
