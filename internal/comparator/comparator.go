package comparator

import (
	"github.com/nikitakuchur/testpoint/internal/io/readers/respreader"
	"github.com/nikitakuchur/testpoint/internal/sender"
	"github.com/nikitakuchur/testpoint/internal/strdiff"
	"log"
	"sync"
)

// RespDiff is the result of comparing two response records.
type RespDiff struct {
	Rec1  respreader.RespRecord
	Rec2  respreader.RespRecord
	Diffs map[string][]strdiff.Diff
}

// Comparator is responsible for performing comparison of two responses.
type Comparator interface {
	Compare(resp1, resp2 sender.Response) (map[string][]strdiff.Diff, error)
}

// CompareResponses compares responses from the given channels using the specified response comparator.
func CompareResponses(records1, records2 <-chan respreader.RespRecord, numComparisons int, comparator Comparator, workers int) <-chan RespDiff {
	responsesToCompare := matchResponses(records1, records2, numComparisons)

	output := make(chan RespDiff)

	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for responses := range responsesToCompare {
				compareRecords(responses[0], responses[1], comparator, output)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

func matchResponses(records1, records2 <-chan respreader.RespRecord, numComparisons int) <-chan []respreader.RespRecord {
	matchedResponses := make(chan []respreader.RespRecord)

	go func() {
		defer close(matchedResponses)

		buffer := make(map[uint64]respreader.RespRecord)

		count := 0
		isRecords1Closed, isRecords2Closed := false, false

		for {
			if numComparisons > 0 && count >= numComparisons {
				// we reached the specified number of comparisons
				return
			}

			if isRecords1Closed && isRecords2Closed {
				break
			}

			select {
			case rec1, ok := <-records1:
				if !ok {
					isRecords1Closed = true
					continue
				}
				rec2, ok := buffer[rec1.ReqHash]
				if !ok {
					// we don't have the second record yet, so we need to put this one aside
					buffer[rec1.ReqHash] = rec1
					break
				}
				delete(buffer, rec1.ReqHash)

				// we have both records, let's send them to compare
				matchedResponses <- []respreader.RespRecord{rec1, rec2}
				count++
			case rec2, ok := <-records2:
				if !ok {
					isRecords2Closed = true
					continue
				}
				rec1, ok := buffer[rec2.ReqHash]
				if !ok {
					buffer[rec2.ReqHash] = rec2
					break
				}
				delete(buffer, rec2.ReqHash)

				matchedResponses <- []respreader.RespRecord{rec1, rec2}
				count++
			}
		}

		for _, rec := range buffer {
			log.Printf("request with hash=%v is missing the second response", rec.ReqHash)
		}
	}()

	return matchedResponses
}

func compareRecords(x, y respreader.RespRecord, comparator Comparator, output chan<- RespDiff) {
	diffs := make(map[string][]strdiff.Diff)
	defer func() {
		if len(diffs) != 0 {
			output <- RespDiff{x, y, diffs}
		}
	}()

	resp1 := sender.Response{Status: x.RespStatus, Body: x.RespBody}
	resp2 := sender.Response{Status: y.RespStatus, Body: y.RespBody}

	respDiffs, err := comparator.Compare(resp1, resp2)
	if err != nil {
		log.Printf("%v, the records with hash=%v were skipped", x.ReqHash, err)
		return
	}

	for k, v := range respDiffs {
		diffs[k] = v
	}
}
