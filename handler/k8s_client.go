package handler

import (
	"fmt"
	"k8s-monitoring-tool/models"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	"log"
	"strconv"
)

type K8sClient struct {
	clientSet *kubernetes.Clientset
}

func NewK8sClient() *K8sClient {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	cset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &K8sClient{clientSet: cset}
}

func (kc *K8sClient) GetAllPods(namespace string) (models.Pods, error) {
	var pods models.Pods

	podList, err := kc.clientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Could not fetch pod list from kubernetes cluster: %s", err)
	}

	podsInCluster := podList.Items
	if len(podsInCluster) > 0 {
		pods = make(models.Pods, len(podsInCluster))
	}

	for _, podInCluster := range podsInCluster {
		podStatus := &models.PodStatus{
			Description: makeStringPtr(podInCluster.Status.Message),
			HostIP:      makeStringPtr(podInCluster.Status.HostIP),
			Phase:       makeStringPtr(string(podInCluster.Status.Phase)),
			PodIP:       makeStringPtr(podInCluster.Status.PodIP),
		}

		pod := &models.Pod{
			Name:   makeStringPtr(podInCluster.Name),
			Status: podStatus,
		}

		pods = append(pods, pod)
	}

	return pods, nil
}

func (kc *K8sClient) UpdatePodDeployment(deployment *models.PodDeployment) (*models.PodDeployment, error) {
	deploymentsClient := kc.clientSet.AppsV1().Deployments(apiv1.NamespaceAll)

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(*deployment.Name, metav1.GetOptions{})
		if getErr != nil {
			return fmt.Errorf("Failed to get latest version of Deployment: %v", getErr)
		}

		replicasInt64, err := strconv.ParseInt(*deployment.Replicas, 10, 32)
		if err != nil {
			return fmt.Errorf("Failed to convert replicas string to int64: %v", err)
		}

		result.Spec.Replicas = int32Ptr(int32(replicasInt64))
		result.Spec.Template.Spec.Containers[0].Image = *deployment.Image

		_, updateErr := deploymentsClient.Update(result)
		return updateErr
	})
	if retryErr != nil {
		return nil, fmt.Errorf("Update failed: %v", retryErr)
	}

	log.Println("Updated deployment...")

	return deployment, nil
}

func makeStringPtr(v string) *string {
	if v != "" {
		return &v
	}
	return nil
}

func int32Ptr(i int32) *int32 { return &i }
