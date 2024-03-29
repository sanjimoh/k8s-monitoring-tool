CURR_SWAGGER_VER := 0.27.0
SWAGGER := $(shell swagger version | grep $(CURR_SWAGGER_VER) | wc -l 2> /dev/null)

build:
	env GOOS=linux CGO_ENABLED=0 go build -o builds/k8s-monitoring-tool cmd/kmt-server/main.go

run: build
	env PORT=30001 builds/k8s-monitoring-tool

swagger-generate-server:
ifeq ($(SWAGGER), 1)
	swagger generate server -t . -f swagger.yml -A kmt
else
	echo "swagger version: $(CURR_SWAGGER_VER) is not available. Please install!"
endif
