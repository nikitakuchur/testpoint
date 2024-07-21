package comparator

import (
	"github.com/nikitakuchur/testpoint/internal/io/readers/respreader"
	"github.com/nikitakuchur/testpoint/internal/sender"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
)

// RespDiff is the result of comparing two response records.
type RespDiff struct {
	Rec1  respreader.RespRecord
	Rec2  respreader.RespRecord
	Diffs map[string][]diffmatchpatch.Diff
}

// Comparator is responsible for performing comparison of two responses.
type Comparator interface {
	Compare(resp1, resp2 sender.Response) (map[string][]diffmatchpatch.Diff, error)
}

// CompareResponses compares responses from the given channels using the specified response comparator.
// You can limit the number of comparisons by using the n param.
func CompareResponses(records1, records2 <-chan respreader.RespRecord, comparator Comparator, numComparisons int) <-chan RespDiff {
	output := make(chan RespDiff)

	go func() {
		defer close(output)

		cache := make(map[uint64]respreader.RespRecord)

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
				rec2, ok := cache[rec1.ReqHash]
				if !ok {
					// we don't have the second record yet, so we need to put this one aside
					cache[rec1.ReqHash] = rec1
					break
				}
				delete(cache, rec1.ReqHash)

				// we have both records, let's compare them
				compareRecords(rec1, rec2, comparator, output)
				count++
			case rec2, ok := <-records2:
				if !ok {
					isRecords2Closed = true
					continue
				}
				rec1, ok := cache[rec2.ReqHash]
				if !ok {
					cache[rec2.ReqHash] = rec2
					break
				}
				delete(cache, rec2.ReqHash)

				compareRecords(rec1, rec2, comparator, output)
				count++
			}
		}

		for _, rec := range cache {
			log.Printf("request with hash=%v is missing the second response", rec.ReqHash)
		}
	}()

	return output
}

func compareRecords(x, y respreader.RespRecord, comparator Comparator, output chan<- RespDiff) {
	diffs := make(map[string][]diffmatchpatch.Diff)
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
