package diff

import (
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"strconv"
	"strings"
)

// Diff finds the differences between two texts line-by-line.
func Diff(text1, text2 string) []diffmatchpatch.Diff {
	dmp := diffmatchpatch.New()

	lineText1, lineText2, lineArray := dmp.DiffLinesToChars(text1, text2)

	// the diffmatchpatch library has a bug:
	// it incorrectly converts an int array into a string, and it leads to broken diffs
	// here's a small workaround:
	lineText1, lineText2 = fixLine(lineText1), fixLine(lineText2)

	diffs := dmp.DiffMain(lineText1, lineText2, false)

	// the CharsToLines function from the library is also incorrect, so I've written my own version
	diffs = charsToLines(diffs, lineArray)

	return diffs
}

func fixLine(line string) string {
	nums := strings.Split(line, ",")
	sb := strings.Builder{}
	for _, num := range nums {
		i, err := strconv.ParseInt(num, 10, 32)
		if err != nil {
			log.Fatalln("failed to convert a string into a number: ", err)
		}
		sb.WriteRune(rune(i))
	}
	return sb.String()
}

func charsToLines(diffs []diffmatchpatch.Diff, lineArray []string) []diffmatchpatch.Diff {
	result := make([]diffmatchpatch.Diff, 0, len(diffs))
	for _, d := range diffs {
		text := make([]string, len(d.Text))

		for i, r := range d.Text {
			text[i] = lineArray[r]
		}

		d.Text = strings.Join(text, "")
		result = append(result, d)
	}
	return result
}
