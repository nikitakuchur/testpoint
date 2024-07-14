package comparator

import (
	"encoding/json"
	"github.com/sergi/go-diff/diffmatchpatch"
	"testpoint/internal/sender"
	jsonutils "testpoint/internal/utils/json"
)

type DefaultComparator struct {
	ignoreOrder bool
}

func NewDefaultComparator(ignoreOrder bool) DefaultComparator {
	return DefaultComparator{ignoreOrder}
}

func (c DefaultComparator) Compare(resp1, resp2 sender.Response) (map[string][]diffmatchpatch.Diff, error) {
	result := make(map[string][]diffmatchpatch.Diff)

	if resp1.Status != resp2.Status {
		result["status"] = diff(resp1.Status, resp2.Status)
		return result, nil
	}

	if resp1.Body != resp2.Body {
		body1, body2 := resp1.Body, resp2.Body
		if json.Valid([]byte(body1)) && json.Valid([]byte(body2)) {
			body1 = jsonutils.ReformatJson(resp1.Body, c.ignoreOrder)
			body2 = jsonutils.ReformatJson(resp2.Body, c.ignoreOrder)
		}
		// we need to check if the bodies are equal again after reformating the JSONs
		if body1 != body2 {
			result["body"] = diff(body1, body2)
		}
	}

	return result, nil
}
