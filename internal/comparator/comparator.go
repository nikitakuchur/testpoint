package comparator

import (
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"testpoint/internal/io/readers/respreader"
	"testpoint/internal/sender"
)

type RespDiff struct {
	Rec1  respreader.RespRecord
	Rec2  respreader.RespRecord
	Diffs map[string][]diffmatchpatch.Diff
}

func CompareResponses(records1, records2 <-chan respreader.RespRecord, respComparator RespComparator) <-chan RespDiff {
	output := make(chan RespDiff)

	go func() {
		cache := make(map[uint64]respreader.RespRecord)

		isRecords1Closed, isRecords2Closed := false, false
		for {
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
				compareRecords(rec1, rec2, respComparator, output)
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

				compareRecords(rec1, rec2, respComparator, output)
			}
		}

		for _, rec := range cache {
			log.Printf("request with hash=%v is missing the second response", rec.ReqHash)
		}

		close(output)
	}()

	return output
}

func compareRecords(rec1, rec2 respreader.RespRecord, respComparator RespComparator, output chan<- RespDiff) {
	diffs := make(map[string][]diffmatchpatch.Diff)
	defer func() {
		if len(diffs) != 0 {
			output <- RespDiff{rec1, rec2, diffs}
		}
	}()

	resp1 := sender.Response{Status: rec1.RespStatus, Body: rec1.RespBody}
	resp2 := sender.Response{Status: rec2.RespStatus, Body: rec2.RespBody}

	respDiffs, err := respComparator(resp1, resp2)
	if err != nil {
		log.Printf("%v, the records with hash=%v were skipped", rec1.ReqHash, err)
		return
	}

	for k, v := range respDiffs {
		diffs[k] = v
	}
}