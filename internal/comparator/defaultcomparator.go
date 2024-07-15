package comparator

import (
	"encoding/json"
	"github.com/sergi/go-diff/diffmatchpatch"
	"testpoint/internal/diff"
	"testpoint/internal/sender"
	jsonutils "testpoint/internal/utils/json"
)

type DefaultComparator struct {
	ignoreOrder bool
}

func NewDefaultComparator(ignoreOrder bool) DefaultComparator {
	return DefaultComparator{ignoreOrder}
}

func (c DefaultComparator) Compare(x, y sender.Response) (map[string][]diffmatchpatch.Diff, error) {
	result := make(map[string][]diffmatchpatch.Diff)

	if x.Status != y.Status {
		result["status"] = diff.Diff(x.Status, y.Status)
		return result, nil
	}

	if x.Body != y.Body {
		body1, body2 := x.Body, y.Body
		if json.Valid([]byte(body1)) && json.Valid([]byte(body2)) {
			body1 = jsonutils.ReformatJson(x.Body, c.ignoreOrder, []string{})
			body2 = jsonutils.ReformatJson(y.Body, c.ignoreOrder, []string{})
		}
		// we need to check if the bodies are equal again after reformating the JSONs
		if body1 != body2 {
			result["body"] = diff.Diff(body1, body2)
		}
	}

	return result, nil
}
