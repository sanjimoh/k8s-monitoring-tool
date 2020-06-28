package controller

import (
	"fmt"
	"k8s-monitoring-tool/configuration"
	"k8s-monitoring-tool/handler"
)

type K8sMonitoringController struct {
	MonitoringHandler *handler.K8sMonitoringHandler
}

func NewK8sMonitoringController() (*K8sMonitoringController, error) {
	config, err := configuration.ParseEnvConfiguration()
	if err != nil {
		return nil, fmt.Errorf("Could not parse k8s monitoring tool service config: %s", err)
	}

	kubernetesClient, err := handler.NewK8sClient(config)
	if err != nil {
		return nil, fmt.Errorf("Creating k8s client failed: %v", err)
	}

	handler, err := handler.NewK8sMonitoringHandler(kubernetesClient)
	if err != nil {
		return nil, fmt.Errorf("Creating k8s monitoring handler failed: %v", err)
	}

	return &K8sMonitoringController{MonitoringHandler: handler}, nil
}
