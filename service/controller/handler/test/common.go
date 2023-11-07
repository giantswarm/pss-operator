package test

import (
	"github.com/giantswarm/microerror"
	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func toCluster(v any) (*capiv1beta1.Cluster, error) {
	c, ok := v.(*capiv1beta1.Cluster)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &capiv1beta1.Cluster{}, v)
	}
	if c == nil {
		return nil, microerror.Maskf(wrongTypeError, "nil cannot be converted to custom resource")
	}
	return c, nil
}
