package reporter

import (
	"fmt"
	"github.com/nikitakuchur/testpoint/internal/comparator"
	"github.com/nikitakuchur/testpoint/internal/io/readers/respreader"
	"github.com/nikitakuchur/testpoint/internal/strdiff"
	"log"
	"strings"
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
		sb.WriteString("req2:\n")
		writeRequest(&sb, d.Rec2)
	}
	sb.WriteString(fmt.Sprintf("\nhash: \t%d\n", d.Rec1.ReqHash))

	for k, v := range d.Diffs {
		sb.WriteString(fmt.Sprintf("\n%s:\n", k))

		t := strdiff.DiffToPrettyText(shortenDiff(v))
		t = strings.ReplaceAll(t, "\n", "\n\t")

		sb.WriteString("\t")
		sb.WriteString(t)
	}

	return sb.String()
}

func writeRequest(sb *strings.Builder, rec respreader.RespRecord) {
	sb.WriteString(fmt.Sprintf("\turl: %s\n", rec.ReqUrl))
	sb.WriteString(fmt.Sprintf("\tmethod: %s\n", rec.ReqMethod))
	if rec.ReqHeaders != "" {
		sb.WriteString(fmt.Sprintf("\theaders: %s\n", rec.ReqHeaders))
	}
	if rec.ReqBody != "" {
		sb.WriteString(fmt.Sprintf("\tbody: %s\n", rec.ReqBody))
	}
}

func shortenDiff(diff []strdiff.Diff) []strdiff.Diff {
	result := make([]strdiff.Diff, len(diff))
	for i, d := range diff {
		result[i].Operation = d.Operation
		if strings.HasSuffix(d.Text, "\n") {
			result[i].Text = d.Text
		} else {
			result[i].Text = d.Text + "\n"
		}

		if d.Operation != strdiff.DiffEqual {
			continue
		}
		substrings := strings.Split(d.Text, "\n")
		if len(substrings) > 8 {
			switch i {
			case 0:
				removedLines := len(substrings) - 4
				tail := strings.Join(substrings[len(substrings)-4:], "\n")
				result[i].Text = fmt.Sprintf("... // %d identical lines\n%s", removedLines, tail)
			case len(diff) - 1:
				removedLines := len(substrings) - 3
				head := strings.Join(substrings[:3], "\n")
				result[i].Text = fmt.Sprintf("%s\n... // %d identical lines\n", head, removedLines)
			default:
				removedLines := len(substrings) - 7
				head := strings.Join(substrings[:3], "\n")
				tail := strings.Join(substrings[len(substrings)-4:], "\n")
				result[i].Text = fmt.Sprintf("%s\n... // %d identical lines\n%s", head, removedLines, tail)
			}
		}
	}
	return result
}
