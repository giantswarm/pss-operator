package pssversion

import (
	"context"

	"github.com/Masterminds/semver/v3"
	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/k8smetadata/pkg/label"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v7/pkg/controller/context/resourcecanceledcontext"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

/*
NOTE: The pssCutoffVersion has to match the ones defined in app-admission-controller.
See https://github.com/giantswarm/app-admission-controller/blob/master/pkg/app/mutate_app_psp_removal.go
*/
var (
	// pssCutoffVersion represents the first & lowest Giant Swarm Release
	// version which does not support PodSecurityPolicies.
	pssCutoffVersion, _ = semver.NewVersion("v19.3.0")
)

const (
	pspLabelKey = "policy.giantswarm.io/psp-status"
	pspLabelVal = "disabled"
)

func (r *Handler) EnsureCreated(ctx context.Context, obj interface{}) error {
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

	releaseSemver, err := semver.NewVersion(releaseVersion)
	if err != nil {
		return microerror.Maskf(executionFailedError, "ReleaseVersion %q is not a valid semver", releaseVersion)
	}

	if releaseSemver.LessThan(pssCutoffVersion) {
		r.logger.Debugf(ctx, "Cluster %q version %q does not require any action", cluster.Name, releaseVersion)
		r.logger.Debugf(ctx, "cancelling resource")
		return nil
	}

	// Label every App belonging to this cluster, forcing them to going throught admission process.
	r.logger.Debugf(ctx, "Cluster %q release version >=%s, adding labels to managed Apps...", cluster.Name, pssCutoffVersion)
	appList := &v1alpha1.AppList{}
	err = r.k8sclient.CtrlClient().List(ctx, appList, &client.ListOptions{Namespace: cluster.Name})
	if err != nil {
		return microerror.Mask(err)
	}

	var patchErrorCount = 0
	for _, app := range appList.Items {
		labelValue, ok := app.Labels[pspLabelKey]
		if ok && labelValue == pspLabelVal {
			continue
		}

		a := app
		a.Labels[pspLabelKey] = pspLabelVal
		err = r.k8sclient.CtrlClient().Update(ctx, &a, &client.UpdateOptions{})
		if err != nil {
			r.logger.Errorf(ctx, err, "error updating App %q for Cluster %q", app.Name, cluster.Name)
			patchErrorCount++
			continue
		}
		r.logger.Debugf(ctx, "added label to App %s/%s", app.Namespace, app.Name)
	}
	r.logger.Debugf(ctx, "finished adding labels for Apps belonging to %q", cluster.Name)

	if patchErrorCount > 0 {
		resourcecanceledcontext.SetCanceled(ctx)
		return microerror.Maskf(executionFailedError, "encountered %d errors while patching apps", patchErrorCount)
	}
	return nil
}