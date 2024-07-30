package main

import (
	"fmt"
	"github.com/nikitakuchur/testpoint/internal/comparator"
	"github.com/nikitakuchur/testpoint/internal/io/readers/respreader"
	"github.com/nikitakuchur/testpoint/internal/io/writers/reporter"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

type compareConfig struct {
	file1          string
	file2          string
	numComparisons int
	comparator     string
	workers        int
	ignoreOrder    bool
	csvReport      string
}

func (c compareConfig) String() string {
	comp := c.comparator
	if comp == "" {
		comp = "default"
	}
	numComparisons := "all"
	if c.numComparisons > 0 {
		numComparisons = strconv.Itoa(c.numComparisons)
	}
	return fmt.Sprintf(
		"file1: %v, file2: %v, numComparisons: %v, comparator: %v, workers: %v, ignoreOrder: %v, csvReport: %v",
		c.file1, c.file2, numComparisons, comp, c.workers, c.ignoreOrder, c.csvReport,
	)
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
			diffs := comparator.CompareResponses(
				records1,
				records2,
				conf.numComparisons,
				createComparator(conf.comparator, conf.ignoreOrder),
				conf.workers,
			)

			reporters := []reporter.Reporter{reporter.NewLogReporter(log.Default())}

			if conf.csvReport != "" {
				reporters = append(reporters, reporter.NewCsvReporter(conf.csvReport))
			}

			reporter.GenerateReport(diffs, reporters...)

			log.Println("completed")
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&conf.numComparisons, "num-comparisons", "n", 0, "number of comparisons to perform")
	flags.StringVarP(&conf.comparator, "comparator", "c", "", "JavaScript file with a response comparator")
	flags.IntVarP(&conf.workers, "workers", "w", 8, "number of workers to compare responses")
	flags.BoolVar(&conf.ignoreOrder, "ignore-order", false, "enable this flag if you want to ignore array order during comparison")
	flags.StringVar(&conf.csvReport, "csv-report", "", "output a comparison report to a CSV file")

	return cmd
}

func createComparator(filepath string, ignoreOrder bool) comparator.Comparator {
	if filepath == "" {
		return comparator.NewDefaultComparator(ignoreOrder)
	}
	script := readComparatorScript(filepath)
	comp, err := comparator.NewScriptComparator(script, ignoreOrder)
	if err != nil {
		log.Fatalln(err)
	}
	return &comp
}

func readComparatorScript(filename string) string {
	script, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln("cannot read the comparator script: ", err)
	}
	return string(script)
}
