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
		return nil, fmt.Errorf("Fetching all pods for the namespace: %s failed: %v", namespace, err)
	}

	return
}

func (kmh *K8sMonitoringHandler) GetV1PodsUnderLoad(namespace string, cpuThreshold string, memoryThreshold string) (pods models.Pods, err error) {
	pods, err = kmh.k8sClient.GetAllPodsUnderLoad(namespace, cpuThreshold, memoryThreshold)
	if err != nil {
		return nil, fmt.Errorf("Fetching pods breaching thresholds for the namespace: %s failed: %v", namespace, err)
	}

	return
}

func (kmh *K8sMonitoringHandler) PutV1Pod(deployment *models.PodDeployment) (pods *models.PodDeployment, err error) {
	pods, err = kmh.k8sClient.UpdatePodDeployment(deployment)
	if err != nil {
		return nil, fmt.Errorf("Update to pod deployment failed: %v", err)
	}

	return
}