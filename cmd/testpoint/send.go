package main

import (
	"fmt"
	"github.com/nikitakuchur/testpoint/internal/filter"
	"github.com/nikitakuchur/testpoint/internal/io/readers/reqreader"
	"github.com/nikitakuchur/testpoint/internal/io/writers/respwriter"
	"github.com/nikitakuchur/testpoint/internal/sender"
	"github.com/nikitakuchur/testpoint/internal/transformer"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

type sendConfig struct {
	input          string
	numRequests    int
	noHeader       bool
	urls           []string
	transformation string
	workers        int
	outputDir      string
}

func (c sendConfig) String() string {
	transformation := c.transformation
	if transformation == "" {
		transformation = "default"
	}
	numRequests := "all"
	if c.numRequests > 0 {
		numRequests = strconv.Itoa(c.numRequests)
	}
	return fmt.Sprintf(
		"input: %v, numRequests: %v, noHeader: %v, urls: %v, transformation: %v, workers: %v, outputDir: %v",
		c.input, numRequests, c.noHeader, c.urls, transformation, c.workers, c.outputDir,
	)
}

func newSendCmd() *cobra.Command {
	var conf sendConfig

	cmd := &cobra.Command{
		Use:   "send [flags] <input> <url>...",
		Short: "Send prepared requests to specified REST endpoints",
		Long:  "Send requests from the given input (CSV file or directory of CSV files) to the specified URLs and collect the responses in output files.",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			conf.input = args[0]
			conf.urls = args[1:]

			log.Printf("configuration: {%v}\n", conf)
			log.Println("starting to process the requests...")

			records := reqreader.ReadRequests(conf.input, !conf.noHeader, conf.numRequests)
			records = filter.Filter(records)
			requests := transformer.TransformRequests(conf.urls, records, createReqTransformation(conf.transformation))
			responses := sender.SendRequests(requests, conf.workers)
			respwriter.WriteResponses(responses, conf.outputDir)

			log.Printf("the result was saved in %v", conf.outputDir)
			log.Println("completed")
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&conf.numRequests, "num-requests", "n", 0, "number of requests to process")
	flags.BoolVar(&conf.noHeader, "no-header", false, "enable this flag if your CSV file has no header")
	flags.StringVarP(&conf.transformation, "transformation", "t", "", "JavaScript file with a request transformation")
	flags.IntVarP(&conf.workers, "workers", "w", 1, "number of workers to send requests")
	flags.StringVar(&conf.outputDir, "output-dir", "./", "directory where the output files need to be saved")

	return cmd
}

func createReqTransformation(filepath string) transformer.ReqTransformation {
	if filepath == "" {
		return transformer.DefaultReqTransformation
	}
	script := readTransformationScript(filepath)
	transformation, err := transformer.NewReqTransformation(script)
	if err != nil {
		log.Fatalln(err)
	}
	return transformation
}

func readTransformationScript(filename string) string {
	script, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln("cannot read the transformation script:", err)
	}
	return string(script)
}
