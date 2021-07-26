package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"k8s-monitoring-tool/configuration"
	"k8s-monitoring-tool/models"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
	"strconv"
	"strings"
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

func NewK8sClient(config *configuration.K8sEnvConfig) (*K8sClient, error) {
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", config.K8sConfig)
	if err != nil {
		panic(err.Error())
	}

	cset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("Could not set k8s cluster configuration: %s", err)
	}

	mc, err := metrics.NewForConfig(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("Could not set k8s metrics cluster configuration: %s", err)
	}

	return &K8sClient{clientSet: cset, metricsClientSet: mc}, nil
}

func (kc *K8sClient) GetAllPods(namespace string) (models.Pods, error) {
	var pods models.Pods
	var containers models.PodContainers
	var labels []string

	podList, err := kc.clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Could not fetch pod list from kubernetes cluster: %s", err)
	}

	podsInCluster := podList.Items
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

		for key, value := range podInCluster.Labels {
			labels = append(labels, key+":"+value)
		}

		pod := &models.Pod{
			Name:       makeStringPtr(podInCluster.Name),
			Namespace:  makeStringPtr(podInCluster.Namespace),
			Status:     podStatus,
			Containers: containers,
			Labels:     strings.Join(labels, ", "),
		}

		pods = append(pods, pod)
	}

	return pods, nil
}

func (kc *K8sClient) UpdatePodDeployment(deployment *models.PodDeployment) (*models.PodDeployment, error) {
	deploymentName := deployment.Name
	replicas := deployment.Replicas
	imageNameVer := deployment.Image
	labelKey := deployment.AffinityKey
	labelValues := deployment.AffinityValues

	if len(deploymentName) > 0 {
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			deploymentsClient := kc.clientSet.AppsV1().Deployments(deployment.Namespace)

			result, getErr := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
			if getErr != nil {
				log.Printf("Error while fetching deployment: %v", getErr)
				return fmt.Errorf("Failed to get latest version of Deployment: %v", getErr)
			}

			// Handling update of number of replicas
			if len(replicas) > 0 {
				replicasInt64, err := strconv.ParseInt(replicas, 10, 32)
				if err != nil {
					return fmt.Errorf("Failed to convert replicas string to int64: %v", err)
				}

				result.Spec.Replicas = int32Ptr(int32(replicasInt64))
			}

			// Handling update of image name & version
			if len(imageNameVer) > 0 {
				result.Spec.Template.Spec.Containers[0].Image = imageNameVer
			}

			// Handling Pod anti-affinity
			if len(labelKey) > 0 && len(labelValues) > 0 {
				affinityTerms := make([]apiv1.WeightedPodAffinityTerm, 1)
				for _, term := range affinityTerms {
					term.Weight = 50
					labelRequirements := make([]metav1.LabelSelectorRequirement, 1)
					for _, label := range labelRequirements {
						label.Key = labelKey
						label.Operator = metav1.LabelSelectorOpIn
						label.Values = strings.Split(labelValues, ",")
					}
					term.PodAffinityTerm.TopologyKey = "kubernetes.io/hostname"
					term.PodAffinityTerm.LabelSelector.MatchExpressions = labelRequirements
				}
				result.Spec.Template.Spec.Affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution = affinityTerms
			}

			_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
			return updateErr
		})
		if retryErr != nil {
			return nil, fmt.Errorf("Update failed: %v", retryErr)
		}

		log.Println("Updated deployment...")
	}

	return deployment, nil
}

func (kc *K8sClient) GetAllPodsUnderLoad(namespace string, cpuThreshold string, memoryThreshold string) (models.Pods, error) {
	var pods models.Pods
	var containers models.PodContainers
	var labels []string

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

				for key, value := range podMetric.Labels {
					labels = append(labels, key+":"+value)
				}

				pod := &models.Pod{
					Name:       makeStringPtr(podMetric.Name),
					Namespace:  makeStringPtr(podMetric.Namespace),
					Status:     podStatus,
					Containers: containers,
					Labels:     strings.Join(labels, ", "),
				}

				pods = append(pods, pod)
			}
		}
	}

	return pods, nil
}

func (kc *K8sClient) GetPodsLog(namespace, podName, containerName string) (string, error) {
	req := kc.clientSet.CoreV1().Pods(namespace).GetLogs(podName, &apiv1.PodLogOptions{Container: containerName})
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return "", errors.New("error in opening stream")
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", errors.New("error in copy information from podLogs to buf")
	}

	return buf.String(), nil
}

func (kc *K8sClient) AlertOnMatchInPodsLog(namespace, podName, containerName string) (*models.PodLogs, error) {
	factory := informers.NewSharedInformerFactory(kc.clientSet, 0)
	informer := factory.Core().V1().Pods().Informer()
	stopper := make(chan struct{})
	defer close(stopper)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(metav1.Object)
			log.Printf("New Pod Added to Store: %s", mObj.GetName())
		},
	})
	informer.Run(stopper)
	return nil, nil
}

func makeStringPtr(v string) *string {
	if v != "" {
		return &v
	}
	return nil
}

func int32Ptr(i int32) *int32 { return &i }
