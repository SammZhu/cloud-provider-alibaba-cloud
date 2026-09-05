package clbv1

import (
	"os"
	"strings"

	"k8s.io/klog/v2"

	"k8s.io/cloud-provider-alibaba-cloud/pkg/model"
)

// DefaultInstanceChargeTypeEnv names the billing model to assume for a Service
// that carries no instance-charge-type annotation.
const DefaultInstanceChargeTypeEnv = "ALIBABA_CLOUD_DEFAULT_INSTANCE_CHARGE_TYPE"

// defaultInstanceChargeType is what such a Service should be treated as wanting.
//
// Upstream answers PayByCLCU unconditionally, because InstanceChargeType's
// IsPayByCLCU() reports true for the empty string.  On a cloud that offers only
// the older PayBySpec model that is not a harmless default: the local model then
// disagrees with every load balancer that exists, so the reconciler calls
// ModifyLoadBalancerInstanceChargeType, and where that API is absent the call
// fails with InvalidAction.NotFound and aborts the entire reconcile.  No Service
// ever reaches an EXTERNAL-IP, and the failure reads like a permissions problem
// rather than a billing-model one.
//
// Unset, this returns PayByCLCU — byte-for-byte upstream behaviour, so public
// cloud is untouched.  A per-Service annotation still overrides whatever this
// returns.
func defaultInstanceChargeType() model.InstanceChargeType {
	v := strings.TrimSpace(os.Getenv(DefaultInstanceChargeTypeEnv))
	switch {
	case v == "":
		return model.PayByCLCU
	case strings.EqualFold(v, string(model.PayBySpec)):
		return model.PayBySpec
	case strings.EqualFold(v, string(model.PayByCLCU)):
		return model.PayByCLCU
	default:
		// More likely a typo than an intention, and falling back in silence
		// would hide it for as long as the cluster runs.
		klog.Warningf("%s=%q is neither PayBySpec nor PayByCLCU; using PayByCLCU",
			DefaultInstanceChargeTypeEnv, v)
		return model.PayByCLCU
	}
}
