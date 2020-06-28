package configuration

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type K8sEnvConfig struct {
	K8sConfig string `envconfig:"KUBE_CONFIG" required:"true"`
}

func ParseEnvConfiguration() (conf *K8sEnvConfig, err error) {
	conf = &K8sEnvConfig{}
	if err := envconfig.Process("", conf); err != nil {
		return nil, fmt.Errorf("Environment variables are not provided: %v ", err)
	}
	return
}
