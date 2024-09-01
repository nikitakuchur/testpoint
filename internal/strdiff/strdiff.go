package strdiff

import (
	"strings"
)

type Operation int8

const (
	// DiffDelete item represents a delete diff
	DiffDelete Operation = -1
	// DiffInsert item represents an insert diff
	DiffInsert Operation = 1
	// DiffEqual item represents an equal diff
	DiffEqual Operation = 0
)

// Diff represents one diff operation
type Diff struct {
	Operation Operation
	Text      string
}

// CalculateLineDiff calculates the differences between two texts line-by-line
func CalculateLineDiff(text1, text2 string) []Diff {
	lineText1, lineText2, lineArray := linesToRunes(text1, text2)
	diffs := CalculateCharDiff(lineText1, lineText2)
	diffs = charsToLines(diffs, lineArray)
	return diffs
}

func linesToRunes(text1, text2 string) ([]rune, []rune, []string) {
	lineToIndex := make(map[string]int)
	var uniqueLines []string

	lines1 := stringToLines(text1, lineToIndex, &uniqueLines)
	lines2 := stringToLines(text2, lineToIndex, &uniqueLines)

	return lines1, lines2, uniqueLines
}

func stringToLines(text string, lineToIndex map[string]int, uniqueLines *[]string) []rune {
	var result []rune

	start := 0
	end := 0
	for end < len(text) {
		if end == len(text)-1 || text[end] == '\n' {
			line := text[start : end+1]
			index, ok := lineToIndex[line]
			if !ok {
				index = len(*uniqueLines)
				*uniqueLines = append(*uniqueLines, line)
			}
			lineToIndex[line] = index
			result = append(result, rune(index))

			end++
			start = end
			continue
		}
		end++
	}

	return result
}

func charsToLines(diffs []Diff, lineArray []string) []Diff {
	for i, d := range diffs {
		text := make([]string, len(d.Text))

		for i, r := range d.Text {
			text[i] = lineArray[r]
		}

		diffs[i].Text = strings.Join(text, "")
	}
	return diffs
}

func CalculateCharDiff(r1, r2 []rune) []Diff {
	// Trim off the common prefix to speed up the algorithm
	commonLength := commonPrefixLength(r1, r2)
	commonPrefix := r1[:commonLength]
	r1 = r1[commonLength:]
	r2 = r2[commonLength:]

	// Trim off the common suffix to speed up the algorithm
	commonLength = commonSuffixLength(r1, r2)
	commonSuffix := r1[len(r1)-commonLength:]
	r1 = r1[:len(r1)-commonLength]
	r2 = r2[:len(r2)-commonLength]

	diff := calculateLcsDiff(r1, r2)

	// Restore the prefix and suffix
	if len(commonPrefix) != 0 {
		diff = append([]Diff{{DiffEqual, string(commonPrefix)}}, diff...)
	}
	if len(commonSuffix) != 0 {
		diff = append(diff, Diff{DiffEqual, string(commonSuffix)})
	}

	return diff
}

func commonPrefixLength(r1, r2 []rune) int {
	length := min(len(r1), len(r2))
	for i := 0; i < length; i++ {
		if r1[i] != r2[i] {
			return i
		}
	}
	return length
}

func commonSuffixLength(r1, r2 []rune) int {
	length := min(len(r1), len(r2))
	for i := 0; i < length; i++ {
		if r1[len(r1)-i-1] != r2[len(r2)-i-1] {
			return i
		}
	}
	return length
}

// calculateLcsDiff calculates difference based on the longest common subsequence
func calculateLcsDiff(r1, r2 []rune) []Diff {
	dp := make([][]int, len(r1)+1)
	for i := range dp {
		dp[i] = make([]int, len(r2)+1)
	}

	for i := 1; i < len(dp); i++ {
		for j := 1; j < len(dp[i]); j++ {
			if r1[i-1] == r2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	i := len(dp) - 1
	j := len(dp[0]) - 1

	var diff []Diff

	var currentOperation Operation
	currentBuilder := strings.Builder{}

	for i > 0 || j > 0 {
		if j > 0 && dp[i][j-1] == dp[i][j] {
			diff = processRune(diff, &currentOperation, &currentBuilder, DiffInsert, r2[j-1])
			j--
		} else if i > 0 && dp[i-1][j] == dp[i][j] {
			diff = processRune(diff, &currentOperation, &currentBuilder, DiffDelete, r1[i-1])
			i--
		} else {
			diff = processRune(diff, &currentOperation, &currentBuilder, DiffEqual, r1[i-1])
			i--
			j--
		}
	}

	if currentBuilder.Len() != 0 {
		diff = append(diff, Diff{currentOperation, reverseString(currentBuilder.String())})
	}

	if len(diff) == 0 {
		return nil
	}
	return reverse(diff)
}

func processRune(diff []Diff, currentOperation *Operation, currentBuilder *strings.Builder, op Operation, r rune) []Diff {
	if *currentOperation == op {
		currentBuilder.WriteRune(r)
	} else {
		if currentBuilder.Len() != 0 {
			diff = append(diff, Diff{*currentOperation, reverseString(currentBuilder.String())})
		}
		*currentOperation = op
		currentBuilder.Reset()
		currentBuilder.WriteRune(r)
	}
	return diff
}

func reverseString(s string) string {
	runes := []rune(s)
	runes = reverse(runes)
	return string(runes)
}

func reverse[T any](arr []T) []T {
	last := len(arr) - 1
	for i := 0; i < len(arr)/2; i++ {
		arr[i], arr[last-i] = arr[last-i], arr[i]
	}
	return arr
}

// DiffToPrettyText converts a []Diff into a colored text report
func DiffToPrettyText(diffs []Diff) string {
	sb := strings.Builder{}

	for _, diff := range diffs {
		text := diff.Text

		switch diff.Operation {
		case DiffInsert:
			sb.WriteString("\x1b[32m")
			sb.WriteString(text)
			sb.WriteString("\x1b[0m")
		case DiffDelete:
			sb.WriteString("\x1b[31m")
			sb.WriteString(text)
			sb.WriteString("\x1b[0m")
		case DiffEqual:
			sb.WriteString(text)
		}
	}

	return sb.String()
}
