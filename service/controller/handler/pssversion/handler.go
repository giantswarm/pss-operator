package pssversion

import (
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

const (
	Name = "pss-version"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
	Provider  string
}

type Handler struct {
	logger    micrologger.Logger
	k8sclient k8sclient.Interface
	provider  string
}

func New(config Config) (*Handler, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	r := &Handler{
		logger:    config.Logger,
		k8sclient: config.K8sClient,
		provider:  config.Provider,
	}

	return r, nil
}

func (r *Handler) Name() string {
	return Name
}
