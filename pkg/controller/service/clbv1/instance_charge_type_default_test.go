package clbv1

import (
	"testing"

	"k8s.io/cloud-provider-alibaba-cloud/pkg/model"
)

func TestDefaultInstanceChargeType(t *testing.T) {
	for _, c := range []struct {
		env  string
		want model.InstanceChargeType
	}{
		// Unset must stay upstream's answer, or this becomes a public-cloud change.
		{"", model.PayByCLCU},
		{"PayBySpec", model.PayBySpec},
		{"paybyspec", model.PayBySpec},
		{"  PayBySpec  ", model.PayBySpec},
		{"PayByCLCU", model.PayByCLCU},
		{"nonsense", model.PayByCLCU},
	} {
		t.Setenv(DefaultInstanceChargeTypeEnv, c.env)
		if got := defaultInstanceChargeType(); got != c.want {
			t.Errorf("%s=%q: got %q, want %q", DefaultInstanceChargeTypeEnv, c.env, got, c.want)
		}
	}
}
