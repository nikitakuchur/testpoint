package main

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"testpoint/internal/comparator"
	"testpoint/internal/io/readers/respreader"
)

type compareConfig struct {
	file1      string
	file2      string
	comparator string
	output     string
}

func (c compareConfig) String() string {
	comp := c.comparator
	if comp == "" {
		comp = "default"
	}
	return fmt.Sprintf("file1: '%v', file2: %v, comparator: %v, output: '%v'", c.file1, c.file2, comp, c.output)
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

			for diff := range diffs {
				printMismatch(diff)
			}

			log.Println("completed")
			//log.Printf("the result is saved in %v", conf.output)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&conf.comparator, "comparator", "c", "", "a JavaScript file with a response comparator")
	flags.StringVarP(&conf.output, "output", "o", "./", "a directory where the output files need to be saved")

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

func printMismatch(d comparator.RespDiff) {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("reqUrl1:\t%s\n", d.Rec1.ReqUrl))
	sb.WriteString(fmt.Sprintf("reqUrl2:\t%s\n", d.Rec2.ReqUrl))
	sb.WriteString(fmt.Sprintf("reqMethod:\t%s\n", d.Rec1.ReqMethod))
	if d.Rec1.ReqHeaders != "" {
		sb.WriteString(fmt.Sprintf("reqHeaders:\t%s\n", d.Rec1.ReqHeaders))
	}
	if d.Rec1.ReqBody != "" {
		sb.WriteString(fmt.Sprintf("reqBody:\t%s\n", d.Rec1.ReqBody))
	}
	sb.WriteString(fmt.Sprintf("reqHash:\t%d\n", d.Rec1.ReqHash))

	dmp := diffmatchpatch.New()
	for k, v := range d.Diffs {
		sb.WriteString(fmt.Sprintf("%s:\n", k))
		sb.WriteString(dmp.DiffPrettyText(v))
	}

	log.Print("MISMATCH:\n", sb.String())
}
