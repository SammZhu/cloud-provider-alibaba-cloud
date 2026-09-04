package slb

// Not a generated file, deliberately.  The struct this extends carries the
// "Changes may cause incorrect behavior and will be lost if the code is
// regenerated" banner, so the behaviour lives beside it rather than inside it.

import (
	"bytes"
	"encoding/json"
)

// UnmarshalJSON accepts both shapes this field is served in.
//
// The documented shape wraps the list in an object:
//
//	"Tags": {"Tag": [{"TagKey": "k", "TagValue": "v"}]}
//
// Apsara Stack gateways return the list directly instead:
//
//	"Tags": [{"TagKey": "k", "TagValue": "v", "InstanceCount": 1}]
//
// Without this, DescribeLoadBalancers fails on any environment where even one
// load balancer carries a tag — the HTTP call succeeds, the body is complete
// and correct, and the client throws it away with
// "readObjectStart: expect { or n, but found [".
//
// The SDK's jsoniter configuration registers a fuzzy extension that bridges
// number/string differences, which is why those cause no trouble; it does not
// bridge object/array, which is why this one does.  jsoniter honours
// json.Unmarshaler, so implementing it here is enough for both decoders.
//
// Accepting the documented shape unchanged means public cloud is unaffected.
func (t *TagsInDescribeLoadBalancers) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimLeft(data, " \t\r\n")

	// Decide on the shape rather than trying one and using the error as a
	// signal: an error here would be indistinguishable from a genuinely
	// malformed body, and silently swallowing those is how a decoder starts
	// returning empty results instead of failures.
	if len(trimmed) > 0 && trimmed[0] == '[' {
		var direct []Tag
		if err := json.Unmarshal(data, &direct); err != nil {
			return err
		}
		t.Tag = direct
		return nil
	}

	// Alias to a type without this method, or the call recurses.
	type plain TagsInDescribeLoadBalancers
	var wrapped plain
	if err := json.Unmarshal(data, &wrapped); err != nil {
		return err
	}
	*t = TagsInDescribeLoadBalancers(wrapped)
	return nil
}
