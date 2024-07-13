package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"testpoint/internal/comparator"
	"testpoint/internal/io/readers/respreader"
	"testpoint/internal/io/writers/reporter"
)

type compareConfig struct {
	file1      string
	file2      string
	comparator string
	csvReport  string
}

func (c compareConfig) String() string {
	comp := c.comparator
	if comp == "" {
		comp = "default"
	}
	return fmt.Sprintf("file1: '%v', file2: %v, comparator: %v, csvReport: '%v'", c.file1, c.file2, comp, c.csvReport)
}

func newCompareCmd() *cobra.Command {
	var conf compareConfig

	cmd := &cobra.Command{
		Use:   "compare [flags] <file1> <file2>",
		Short: "Compare responses and generate a report",
		Long:  "Compare the responses from the given CSV files (the output of the send command) and generate a report on mismatches.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			conf.file1 = args[0]
			conf.file2 = args[1]

			log.Printf("configuration: {%v}\n", conf)
			log.Println("starting to compare the responses...")

			records1 := respreader.ReadResponses(conf.file1)
			records2 := respreader.ReadResponses(conf.file2)
			diffs := comparator.CompareResponses(records1, records2, createRespComparator(conf.comparator))

			reporters := []reporter.Reporter{reporter.NewLogReporter(log.Default())}

			if conf.csvReport != "" {
				reporters = append(reporters, reporter.NewCsvReporter(conf.csvReport))
			}

			reporter.GenerateReport(diffs, reporters...)

			log.Println("completed")
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&conf.comparator, "comparator", "c", "", "a JavaScript file with a response comparator")
	flags.StringVar(&conf.csvReport, "csv-report", "", "output a comparison report to a CSV file")

	return cmd
}

func createRespComparator(filepath string) comparator.RespComparator {
	if filepath == "" {
		return comparator.DefaultRespComparator
	}
	script := readComparatorScript(filepath)
	respComparator, err := comparator.NewRespComparator(script)
	if err != nil {
		log.Fatalln(err)
	}
	return respComparator
}

func readComparatorScript(filename string) string {
	script, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln("cannot read the comparator script:", err)
	}
	return string(script)
}
