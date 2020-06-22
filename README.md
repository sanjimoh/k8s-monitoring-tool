# Kubernetes Monitoring Tool (KMT)
## Scope
Kubernetes monitoring tool is meant for monitoring kubernetes resources in a cluster. It is still work in progress.

The existing codebase - (more to come..)
* Exposes a rest endpoint to it's client using which pod details including its current status could be retrieved from
  a cluster.

## Exposed rest endpoints and payload details
Existing codebase exposes the following endpoints -
* /api/kmt/v1/pods
  This endpoint will fetch all the pods & its status across all the k8s namespaces. 
  
* /api/kmt/v1/pods?namespace=<provide_k8s_namespace_here>
  This endpoint will fetch all the pods for the given k8s namespaces
  
In both the above scenario, the returned response would look like something below -
    
  ```
    [
      {
        "name": "pod-1",
        "status": {
          "phase": "Running",
          "description": "Pod is running",
          "podIp": "192.1.1.1",
          "hostIp": "string"
        }
      },
      {
        "name": "pod-2",
        "status": {
          "phase": "Pending",
          "description": "Pending due to lack of resources",
          "podIp": "",
          "hostIp": ""
        }
      },
      ...
    ]
  ```
  
## Technology Selection
* [Golang](https://golang.org/) used for implementation.
* [Go-Swagger](https://github.com/go-swagger/go-swagger) is used for rest service Swaggerization.
* [client-go](https://github.com/kubernetes/client-go) is used for talking to kubernetes cluster.