package slb

import (
	"errors"
	"fmt"
	"testing"

	sdkerrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
)

// The bodies below are what the gateways actually returned on ste2 on
// 2026-09-13, not invented shapes: an Apsara Stack SLB endpoint has no
// centralised tag service, so ListTagResources is answered with
// InvalidAction.NotFound while DescribeTags works.
func TestIsActionNotFound(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "apsara gateway has no ListTagResources",
			err: sdkerrors.NewServerError(404,
				`{"Code":"InvalidAction.NotFound","Message":"Specified api is not found, please check your url and method.","RequestId":"2D266259-672B-4030-ACA4-48B0F36F6030"}`,
				""),
			want: true,
		},
		{
			name: "a real API error must not trigger the fallback",
			err: sdkerrors.NewServerError(400,
				`{"Code":"InvalidParameter","Message":"The specified parameter is not valid.","RequestId":"x"}`,
				""),
			want: false,
		},
		{
			name: "throttling must not trigger the fallback",
			err: sdkerrors.NewServerError(429,
				`{"Code":"Throttling","Message":"Request was denied due to request throttling.","RequestId":"x"}`,
				""),
			want: false,
		},
		{name: "a plain error is not a server error", err: errors.New("dial tcp: i/o timeout"), want: false},
		{name: "a wrapped server error is still found", err: fmt.Errorf("apply model error: %w",
			sdkerrors.NewServerError(404, `{"Code":"InvalidAction.NotFound","Message":"Specified api is not found","RequestId":"x"}`, "")), want: true},
		{name: "no error", err: nil, want: false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isActionNotFound(c.err); got != c.want {
				t.Fatalf("isActionNotFound(%v) = %v, want %v", c.err, got, c.want)
			}
		})
	}
}
