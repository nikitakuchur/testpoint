package comparator

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"log"
	"strings"
	"testpoint/internal/io/readers/respreader"
)

type Diff struct {
	A     respreader.RespRecord
	B     respreader.RespRecord
	Diffs map[string]string
}

func (d Diff) String() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("reqUrlA:\t%s\n", d.A.ReqUrl))
	sb.WriteString(fmt.Sprintf("reqUrlB:\t%s\n", d.B.ReqUrl))
	sb.WriteString(fmt.Sprintf("reqMethod:\t%s\n", d.A.ReqMethod))
	if d.A.ReqHeaders != "" {
		sb.WriteString(fmt.Sprintf("reqHeaders:\t%s\n", d.A.ReqHeaders))
	}
	if d.A.ReqBody != "" {
		sb.WriteString(fmt.Sprintf("reqBody:\t%s\n", d.A.ReqBody))
	}
	sb.WriteString(fmt.Sprintf("reqHash:\t%d\n", d.A.ReqHash))

	for k, v := range d.Diffs {
		sb.WriteString(fmt.Sprintf("%s:\n", k))
		sb.WriteString(v)
	}

	return sb.String()
}

func CompareResponses(recordsA, recordsB <-chan respreader.RespRecord, respComparator RespComparator) <-chan Diff {
	output := make(chan Diff)

	go func() {
		cache := make(map[uint64]respreader.RespRecord)

		isAClosed, isBClosed := false, false
		for {
			if isAClosed && isBClosed {
				break
			}

			select {
			case a, ok := <-recordsA:
				if !ok {
					isAClosed = true
					continue
				}
				b, ok := cache[a.ReqHash]
				if !ok {
					// we don't have the second record yet, so we need to put this record aside
					cache[a.ReqHash] = a
					break
				}
				delete(cache, a.ReqHash)

				// we have both records, let's compareRecords them
				compareRecords(a, b, respComparator, output)
			case b, ok := <-recordsB:
				if !ok {
					isBClosed = true
					continue
				}
				a, ok := cache[b.ReqHash]
				if !ok {
					cache[b.ReqHash] = b
					break
				}
				delete(cache, b.ReqHash)

				compareRecords(a, b, respComparator, output)
			}
		}
		close(output)

		for _, v := range cache {
			log.Printf("%v: missing second response", v)
		}
	}()

	return output
}

func compareRecords(a, b respreader.RespRecord, respComparator RespComparator, output chan<- Diff) {
	diffs := make(map[string]string)
	defer func() {
		if len(diffs) != 0 {
			output <- Diff{a, b, diffs}
		}
	}()

	if diff := cmp.Diff(a.RespStatus, b.RespStatus); diff != "" {
		diffs["status"] = diff
		return
	}

	// if one of the bodies is not a JSON, then we compareRecords them as strings
	if !json.Valid([]byte(a.RespBody)) || !json.Valid([]byte(b.RespBody)) {
		if diff := cmp.Diff(a.RespBody, b.RespBody); diff != "" {
			diffs["body"] = diff
		}
		return
	}

	respDiffs := respComparator(a, b)
	for k, v := range respDiffs {
		diffs[k] = v
	}
}
