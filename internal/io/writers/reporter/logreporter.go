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

		t := dmp.DiffPrettyText(shortenDiff(v))
		t = strings.ReplaceAll(t, "\n", "\n\t")

		sb.WriteString("\t")
		sb.WriteString(t)
		sb.WriteString("\n")
	}

	return sb.String()
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
