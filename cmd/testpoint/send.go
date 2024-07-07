package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"testpoint/internal/filter"
	"testpoint/internal/io/readers/reqreader"
	"testpoint/internal/io/writers/respwriter"
	"testpoint/internal/sender"
	"testpoint/internal/transformer"
)

type sendConfig struct {
	input          string
	noHeader       bool
	urls           []string
	transformation string
	workers        int
	output         string
}

func (c sendConfig) String() string {
	transformation := c.transformation
	if transformation == "" {
		transformation = "default"
	}
	return fmt.Sprintf(
		"input: '%v', noHeader: %v, urls: %v, transform: %v, workers: %v, output: '%v'",
		c.input, c.noHeader, c.urls, transformation, c.workers, c.output,
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

			records := reqreader.ReadRequests(conf.input, !conf.noHeader)
			records = filter.Filter(records)
			requests := transformer.TransformRequests(conf.urls, records, createTransformation(conf.transformation))
			responses := sender.SendRequests(requests, conf.workers)
			respwriter.WriteResponses(responses, conf.output)

			log.Println("completed")
			log.Printf("the result was saved in %v", conf.output)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&conf.noHeader, "no-header", false, "enable this flag if your CSV file has no header")
	flags.StringVarP(&conf.transformation, "transformation", "t", "", "a JavaScript file with a request transformation")
	flags.IntVarP(&conf.workers, "workers", "w", 1, "a number of workers to send requests")
	flags.StringVarP(&conf.output, "output", "o", "./", "a directory where the output files need to be saved")

	return cmd
}

func createTransformation(filepath string) transformer.Transformation {
	if filepath == "" {
		return transformer.DefaultTransformation
	}
	script := readTransformationScript(filepath)
	transformation, err := transformer.NewTransformation(script)
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
