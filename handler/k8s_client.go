package handler

import (
	"fmt"
	"k8s-monitoring-tool/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

func makeStringPtr(v string) *string {
	if v != "" {
		return &v
	}
	return nil
}
