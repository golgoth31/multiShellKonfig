package namespace

import (
	"errors"

	"github.com/golgoth31/multiShellKonfig/pkg/konfig"
	"github.com/rs/zerolog/log"
)

var (
	errNoReqID    = errors.New("Request ID not set")
	errKubeConfig = errors.New("context not set")
)

func New(curKubeConfig string) (Namespace, error) {
	if curKubeConfig == "" {
		return Namespace{}, errKubeConfig
	}

	return Namespace{
		CurKubeConfig: curKubeConfig,
	}, nil
}

func (ns *Namespace) GetNsList() ([]string, error) {
	log.Debug().Msgf("found config: %s", ns.CurKubeConfig)

	namespaceList, err := konfig.GetNSList(ns.CurKubeConfig)
	if err != nil {
		return namespaceList, err
	}

	return namespaceList, nil
}
