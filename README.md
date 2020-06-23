# Kubernetes Monitoring Tool (KMT)
## Scope
Kubernetes monitoring tool is meant for monitoring kubernetes resources in a cluster. It is still work in progress.

The existing codebase - (more to come..)
* Exposes a rest endpoint to it's client using which pod details including its current status could be retrieved from
  a cluster. Now, the response also includes a list of containers (only name) consisting in the pod.
* Now it is also possible to update a pod deployment through rest endpoints but support is limited to updating of pod 
  deployment number of replicas and pod deployment container image name & version. Other update will be soon introduced.
* Now it is also possible to fetch a list of pods and it's containers which is breaching a passed CPU & memory 
  thresholds.
* Now it is also possible to update pod deployment with pod anti-affinity rules which would ensure pods are running in 
  different Kubernetes nodes. This is implemented based on Kubernetes label selectors. The pod anti-affinity rule 
  currently supported is "preferredDuringSchedulingIgnoredDuringExecution"

## Exposed rest endpoints and payload details
Existing codebase exposes the following endpoints -
* GET /api/kmt/v1/pods
  This endpoint will fetch all the pods & its status across all the k8s namespaces. 
  
* GET /api/kmt/v1/pods?namespace=<provide_k8s_namespace_here>
  This endpoint will fetch all the pods for the given k8s namespaces

* GET /api/kmt/v1/pods?namespace=<provide_k8s_namespace_here>&cpuThreshold=<cpu_threshold>&memoryThreshold=<memory_threshold>
  This endpoint will fetch all the pods for the given k8s namespaces which breaches the given CPU & memory thresholds.

In both the above scenario, the returned response would look like something below -
    
  ```
    [
      {
        "name": "pod-1",
        "namespace": "databricks",
        "labels": "app:pod-1",
        "status": {
          "phase": "Running",
          "description": "Pod is running",
          "podIp": "192.1.1.1",
          "hostIp": "192.1.1.10"
        },
        "containers": [
           {
             "name": "container-0",
             "currentCpuUsage": "",
             "currentMemoryUsage": ""
           }
        ]
      },
      {
        "name": "pod-2",
        "namespace": "databricks",
        "labels": "app:pod-2",
        "status": {
          "phase": "Pending",
          "description": "Pending due to lack of resources",
          "podIp": "",
          "hostIp": ""
        },
        "containers": [
           {
             "name": "container-1",
             "currentCpuUsage": "",
             "currentMemoryUsage": ""
           }
        ]
      },
      ...
    ]
  ```
* PUT /api/kmt/v1/pod
  This endpoint will update a given pod deployment in a kubernetes cluster. Current support is limited to updation of 
  number of pod replicas and image name & version.
  Use "affinityKey" & "affinityValues" if you intend to ensure pod anti-affinity rules. In this case the pods matching
  to "affinityKey" & "affinityValues" would be scheduled in different Kubernetes nodes.
  ```
    PUT /api/kmt/v1/pod
   ```

    Sample request body will be:

    ```
    {
      "name": "apache-cassandra",
      "replicas": "5",
      "image": "ccas-apache:2.5",
      "affinityKey": "app",
      "affinityValues": "ccas-apache,mariadb"
    }
  ```

## Technology Selection
* [Golang](https://golang.org/) used for implementation.
* [Go-Swagger](https://github.com/go-swagger/go-swagger) is used for rest service Swaggerization.
* [client-go](https://github.com/kubernetes/client-go) is used for talking to kubernetes cluster.