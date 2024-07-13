package reporter

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"strings"
	"testpoint/internal/comparator"
	"testpoint/internal/io/readers/respreader"
)

// LogReporter represents a reporter that simply logs all mismatches.
type LogReporter struct {
	logger *log.Logger
}

func NewLogReporter(logger *log.Logger) LogReporter {
	return LogReporter{logger: logger}
}

func (r LogReporter) Report(input <-chan comparator.RespDiff) {
	for diff := range input {
		r.logger.Print("MISMATCH:\n", buildMismatch(diff))
	}
}

func buildMismatch(d comparator.RespDiff) string {
	sb := strings.Builder{}

	if d.Rec1 == d.Rec2 {
		sb.WriteString("req:\n")
		writeRequest(&sb, d.Rec1)
	} else {
		sb.WriteString("req1:\n")
		writeRequest(&sb, d.Rec1)
		sb.WriteString("\nreq2:\n")
		writeRequest(&sb, d.Rec2)
	}

	sb.WriteString(fmt.Sprintf("\nhash: \t%d\n", d.Rec1.ReqHash))

	dmp := diffmatchpatch.New()
	for k, v := range d.Diffs {
		sb.WriteString(fmt.Sprintf("%s:\n", k))

		t := dmp.DiffPrettyText(shortenDiff(v))
		t = strings.ReplaceAll(t, "\n", "\n\t")

		sb.WriteString("\t")
		sb.WriteString(t)
		sb.WriteString("\n")
	}

	return sb.String()
}

func writeRequest(sb *strings.Builder, rec respreader.RespRecord) {
	sb.WriteString(fmt.Sprintf("\turl:\t%s\n", rec.ReqUrl))
	sb.WriteString(fmt.Sprintf("\tmethod:\t%s\n", rec.ReqMethod))
	if rec.ReqHeaders != "" {
		sb.WriteString(fmt.Sprintf("\theaders:\t%s\n", rec.ReqHeaders))
	}
	if rec.ReqBody != "" {
		sb.WriteString(fmt.Sprintf("\tbody:\t%s\n", rec.ReqBody))
	}
}

func shortenDiff(diff []diffmatchpatch.Diff) []diffmatchpatch.Diff {
	result := make([]diffmatchpatch.Diff, len(diff))
	for i, d := range diff {
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
				result[i].Text = fmt.Sprintf("... // %d identical lines\n%s", removedLines, tail)
			case len(diff) - 1:
				removedLines := len(substrings) - 3
				head := strings.Join(substrings[:3], "\n")
				result[i].Text = fmt.Sprintf("%s\n... // %d identical lines", head, removedLines)
			default:
				removedLines := len(substrings) - 6
				head := strings.Join(substrings[:3], "\n")
				tail := strings.Join(substrings[len(substrings)-3:], "\n")
				result[i].Text = fmt.Sprintf("%s\n... // %d identical lines\n%s", head, removedLines, tail)
			}
		}
	}
	return result
}
