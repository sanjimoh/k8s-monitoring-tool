package handler

import (
	"context"
	"fmt"
	"k8s-monitoring-tool/models"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
	"strconv"
)

const (
	//in cores (3 vCPUs)
	cpuThreshold = 3

	//in bytes (1 GB)
	memoryThreshold = 1 * 1024 * 1024 * 1024
)

type K8sClient struct {
	clientSet        *kubernetes.Clientset
	metricsClientSet *metrics.Clientset
}

func NewK8sClient() (*K8sClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("Could not fetch k8s cluster configuration: %s", err)
	}

	cset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Could not set k8s cluster configuration: %s", err)
	}

	mc, err := metrics.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Could not set k8s metrics cluster configuration: %s", err)
	}

	return &K8sClient{clientSet: cset, metricsClientSet: mc}, nil
}

func (kc *K8sClient) GetAllPods(namespace string) (models.Pods, error) {
	var pods models.Pods
	var containers models.PodContainers

	podList, err := kc.clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Could not fetch pod list from kubernetes cluster: %s", err)
	}

	podsInCluster := podList.Items
	if len(podsInCluster) > 0 {
		pods = make(models.Pods, len(podsInCluster))
	}

	for _, podInCluster := range podsInCluster {
		for _, container := range podInCluster.Spec.Containers {
			podContainer := &models.PodContainer{
				CurrentCPUUsage:    nil,
				CurrentMemoryUsage: nil,
				Name:               makeStringPtr(container.Name),
			}
			containers = append(containers, podContainer)
		}

		podStatus := &models.PodStatus{
			Description: makeStringPtr(podInCluster.Status.Message),
			HostIP:      makeStringPtr(podInCluster.Status.HostIP),
			Phase:       makeStringPtr(string(podInCluster.Status.Phase)),
			PodIP:       makeStringPtr(podInCluster.Status.PodIP),
		}

		pod := &models.Pod{
			Name:       makeStringPtr(podInCluster.Name),
			Namespace:  makeStringPtr(podInCluster.Namespace),
			Status:     podStatus,
			Containers: containers,
		}

		pods = append(pods, pod)
	}

	return pods, nil
}

func (kc *K8sClient) UpdatePodDeployment(deployment *models.PodDeployment) (*models.PodDeployment, error) {
	deploymentsClient := kc.clientSet.AppsV1().Deployments(apiv1.NamespaceAll)

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(context.TODO(), *deployment.Name, metav1.GetOptions{})
		if getErr != nil {
			return fmt.Errorf("Failed to get latest version of Deployment: %v", getErr)
		}

		replicasInt64, err := strconv.ParseInt(*deployment.Replicas, 10, 32)
		if err != nil {
			return fmt.Errorf("Failed to convert replicas string to int64: %v", err)
		}

		result.Spec.Replicas = int32Ptr(int32(replicasInt64))
		result.Spec.Template.Spec.Containers[0].Image = *deployment.Image

		_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		return nil, fmt.Errorf("Update failed: %v", retryErr)
	}

	log.Println("Updated deployment...")

	return deployment, nil
}

func (kc *K8sClient) GetAllPodsUnderLoad(namespace string, cpuThreshold string, memoryThreshold string) (models.Pods, error) {
	var pods models.Pods
	var containers models.PodContainers

	cpuThresholdInInt64, err := strconv.ParseInt(cpuThreshold, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert replicas string to int64: %v", err)
	}

	memoryThresholdInInt64, err := strconv.ParseInt(memoryThreshold, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert replicas string to int64: %v", err)
	}

	podMetrics, err := kc.metricsClientSet.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch pod metrics: %v", err)
	}

	for _, podMetric := range podMetrics.Items {
		podContainers := podMetric.Containers
		for _, container := range podContainers {
			cpuQuantity, _ := container.Usage.Cpu().AsInt64()
			memQuantity, _ := container.Usage.Memory().AsInt64()

			// Check if threshold breached
			if cpuQuantity > cpuThresholdInInt64 || memQuantity > memoryThresholdInInt64 {
				podStatus := &models.PodStatus{}

				podContainer := &models.PodContainer{
					CurrentCPUUsage:    makeStringPtr(strconv.FormatInt(cpuQuantity, 10)),
					CurrentMemoryUsage: makeStringPtr(strconv.FormatInt(memQuantity, 10)),
					Name:               makeStringPtr(container.Name),
				}
				containers = append(containers, podContainer)

				pod := &models.Pod{
					Name:       makeStringPtr(podMetric.Name),
					Namespace:  makeStringPtr(podMetric.Namespace),
					Status:     podStatus,
					Containers: containers,
				}

				pods = append(pods, pod)
			}
		}
	}

	return pods, nil
}

func makeStringPtr(v string) *string {
	if v != "" {
		return &v
	}
	return nil
}

func int32Ptr(i int32) *int32 { return &i }
