package controller

import (
	"fmt"
	"k8s-monitoring-tool/handler"
)

type K8sMonitoringController struct {
	MonitoringHandler *handler.K8sMonitoringHandler
}

func NewK8sMonitoringController() (*K8sMonitoringController, error) {
	kubernetesClient := handler.NewK8sClient()

	handler, err := handler.NewK8sMonitoringHandler(kubernetesClient)
	if err != nil {
		return nil, fmt.Errorf("Creating k8s monitoring handler failed: %s", err)
	}

	return &K8sMonitoringController{MonitoringHandler: handler}, nil
}
