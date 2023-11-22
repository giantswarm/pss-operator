package controller

import (
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v7/pkg/controller"
	"github.com/giantswarm/operatorkit/v7/pkg/resource"
	"github.com/giantswarm/operatorkit/v7/pkg/resource/wrapper/metricsresource"
	"github.com/giantswarm/operatorkit/v7/pkg/resource/wrapper/retryresource"
	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/pss-operator/pkg/project"
	"github.com/giantswarm/pss-operator/service/controller/handler/pssversion"
)

type PSSVersionConfig struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
	Provider  string
}

type PSSVersion struct {
	*controller.Controller
}

func NewPSSVersion(config PSSVersionConfig) (*PSSVersion, error) {
	var err error

	handlers, err := newPSSVersionHandlers(config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var operatorkitController *controller.Controller
	{
		c := controller.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
			NewRuntimeObjectFunc: func() client.Object {
				return new(capiv1beta1.Cluster)
			},
			Resources: handlers,
			Name:      project.Name() + "-pss-version-controller",
		}

		operatorkitController, err = controller.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	c := &PSSVersion{
		Controller: operatorkitController,
	}

	return c, nil
}

func newPSSVersionHandlers(config PSSVersionConfig) ([]resource.Interface, error) {
	var err error

	var pssversionResource resource.Interface
	{
		c := pssversion.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		pssversionResource, err = pssversion.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	handlers := []resource.Interface{
		pssversionResource,
	}

	{
		c := retryresource.WrapConfig{
			Logger: config.Logger,
		}

		handlers, err = retryresource.Wrap(handlers, c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	{
		c := metricsresource.WrapConfig{}

		handlers, err = metricsresource.Wrap(handlers, c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return handlers, nil
}
