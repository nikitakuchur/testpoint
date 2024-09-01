package comparator

import (
	"encoding/json"
	"github.com/nikitakuchur/testpoint/internal/sender"
	"github.com/nikitakuchur/testpoint/internal/strdiff"
	jsonutils "github.com/nikitakuchur/testpoint/internal/utils/json"
)

type DefaultComparator struct {
	ignoreOrder bool
}

func NewDefaultComparator(ignoreOrder bool) DefaultComparator {
	return DefaultComparator{ignoreOrder}
}

func (c DefaultComparator) Compare(x, y sender.Response) (map[string][]strdiff.Diff, error) {
	result := make(map[string][]strdiff.Diff)

	if x.Status != y.Status {
		result["status"] = strdiff.CalculateLineDiff(x.Status, y.Status)
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
			result["body"] = strdiff.CalculateLineDiff(body1, body2)
		}
	}

	return result, nil
}
