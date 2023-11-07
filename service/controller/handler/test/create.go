package test

import (
	"context"

	"github.com/giantswarm/k8smetadata/pkg/label"
	"github.com/giantswarm/microerror"
)

func (r *Handler) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "KUBA: EnsureCreated: %+v", obj)
	cluster, err := toCluster(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	releaseVersion, ok := cluster.Labels[label.ReleaseVersion]
	if !ok {
		r.logger.Debugf(ctx, "could not determine release version because Cluster %q/%q does not have a %q label",
			cluster.Namespace, cluster.Name, label.ReleaseVersion)
		r.logger.Debugf(ctx, "cancelling resource")
		return nil
	}

	// TODO: depending on release version ensure empty or filled ConfigMap

	return nil
}
