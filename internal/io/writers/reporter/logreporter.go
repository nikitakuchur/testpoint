package reporter

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"strings"
	"testpoint/internal/comparator"
)

// LogReporter represents a reporter that simply logs all mismatches.
type LogReporter struct {
}

func (r LogReporter) report(input <-chan comparator.RespDiff) {
	for diff := range input {
		printMismatch(diff)
	}
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
		sb.WriteString(fmt.Sprintf("%s:", k))
		sb.WriteString(dmp.DiffPrettyText(shortenDiffs(v)))
	}

	log.Print("MISMATCH:\n", sb.String())
}

func shortenDiffs(diffs []diffmatchpatch.Diff) []diffmatchpatch.Diff {
	result := make([]diffmatchpatch.Diff, len(diffs))
	for i, d := range diffs {
		result[i].Type = d.Type
		result[i].Text = d.Text
		if d.Type != diffmatchpatch.DiffEqual {
			continue
		}
		substrings := strings.Split(d.Text, "\n")
		if len(substrings) > 8 {
			switch i {
			case 0:
				removedLines := len(substrings) - 3
				tail := strings.Join(substrings[len(substrings)-3:], "\n")
				result[i].Text = fmt.Sprintf("\n... // %d identical lines\n %s", removedLines, tail)
			case len(diffs) - 1:
				removedLines := len(substrings) - 3
				head := strings.Join(substrings[:3], "\n")
				result[i].Text = fmt.Sprintf(" %s \n... // %d identical lines\n", head, removedLines)
			default:
				removedLines := len(substrings) - 6
				head := strings.Join(substrings[:3], "\n")
				tail := strings.Join(substrings[len(substrings)-3:], "\n")
				result[i].Text = fmt.Sprintf(" %s \n... // %d identical lines\n %s", head, removedLines, tail)
			}
		}
	}
	return result
}
