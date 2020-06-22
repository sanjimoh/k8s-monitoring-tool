// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/go-openapi/swag"
	"k8s-monitoring-tool/controller"
	"k8s-monitoring-tool/models"
	"log"
	"net/http"
	"os"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"k8s-monitoring-tool/restapi/operations"
	"k8s-monitoring-tool/restapi/operations/k8s_monitoring_tool"
)

//go:generate swagger generate server --target ../../k8s-monitoring-tool --name Kmt --spec ../swagger.yml --exclude-main

var KMC *controller.K8sMonitoringController

func configureFlags(api *operations.KmtAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.KmtAPI) http.Handler {
	var err error

	// configure the api here
	api.ServeError = errors.ServeError
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	if KMC == nil {
		KMC, err = controller.NewK8sMonitoringController()
		if err != nil {
			log.Printf("Failed to initialize k8s monitoring tool service, error: %s", err)
			os.Exit(1)
		}
	}

	// API handling for get all pod status & pods which are breaching cpu/memory thresholds
	api.K8sMonitoringToolGetV1PodsHandler = k8s_monitoring_tool.GetV1PodsHandlerFunc(func(params k8s_monitoring_tool.GetV1PodsParams) middleware.Responder {
		var pods models.Pods
		var err error

		namespace, cpuLimit, memoryLimit := validateAndSetParams(params)

		if len(cpuLimit) == 0 && len(memoryLimit) == 0 {
			pods, err = KMC.MonitoringHandler.GetV1Pods(namespace)
			if err != nil {
				return k8s_monitoring_tool.NewGetV1PodsInternalServerError().WithPayload(&models.Error{
					Code:    swag.Int64(500),
					Message: swag.String(err.Error()),
				})
			}
		} else {
			pods, err = KMC.MonitoringHandler.GetV1PodsUnderLoad(namespace, cpuLimit, memoryLimit)
			if err != nil {
				return k8s_monitoring_tool.NewGetV1PodsInternalServerError().WithPayload(&models.Error{
					Code:    swag.Int64(500),
					Message: swag.String(err.Error()),
				})
			}
		}

		return k8s_monitoring_tool.NewGetV1PodsOK().WithPayload(pods)
	})

	api.K8sMonitoringToolPutV1PodHandler = k8s_monitoring_tool.PutV1PodHandlerFunc(func(params k8s_monitoring_tool.PutV1PodParams) middleware.Responder {
		podDeployment, err := KMC.MonitoringHandler.PutV1Pod(params.PodDeployment)
		if err != nil {
			return k8s_monitoring_tool.NewPutV1PodInternalServerError().WithPayload(&models.Error{
				Code:    swag.Int64(500),
				Message: swag.String(err.Error()),
			})
		}

		return k8s_monitoring_tool.NewPutV1PodOK().WithPayload(podDeployment)
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

func validateAndSetParams(params k8s_monitoring_tool.GetV1PodsParams) (string, string, string) {
	var namespace, cpuLimit, memoryLimit string

	//Check passed namespace
	ns := params.Namespace
	if ns == nil {
		namespace = ""
	} else {
		namespace = *params.Namespace
	}

	//Check passed cpu threshold value
	cpuThreshold := params.CPUThreshold
	if cpuThreshold == nil {
		cpuLimit = ""
	} else {
		cpuLimit = *cpuThreshold
	}

	//Check passed memory threshold value
	memoryThreshold := params.MemoryThreshold
	if memoryThreshold == nil {
		memoryLimit = ""
	} else {
		memoryLimit = *memoryThreshold
	}
	return namespace, cpuLimit, memoryLimit
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
