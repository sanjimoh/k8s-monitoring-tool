package handler

import (
	"fmt"
	"k8s-monitoring-tool/models"
)

type K8sMonitoringHandler struct {
	k8sClient *K8sClient
}

func NewK8sMonitoringHandler(kubernetesClient *K8sClient) (*K8sMonitoringHandler, error) {
	return &K8sMonitoringHandler{k8sClient: kubernetesClient}, nil
}

func (kmh *K8sMonitoringHandler) GetV1Pods(namespace string) (pods models.Pods, err error) {
	pods, err = kmh.k8sClient.GetAllPods(namespace)
	if err != nil {
		return nil, fmt.Errorf("Creating k8s monitoring handler failed: %s", err)
	}

	return
}
