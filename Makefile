CURR_SWAGGER_VER := 0.24.0
SWAGGER := $(shell swagger version | grep $(CURR_SWAGGER_VER) | wc -l 2> /dev/null)

build:
	env GOOS=linux CGO_ENABLED=0 go build -o builds/k8s-monitoring-tool cmd/kmt-server/main.go

run:
	env GO111MODULE=on go build -mod=vendor
	env APP_CREDENTIAL_PROVIDER_URL=http://sco-service-mock:1080 ENTITY_NAME=cassandra CONCURRENT_MESSAGE_PROCESSORS=16 PROFILE_SERVICE_ENABLED=false SERVICE_DEBUG_HOST=localhost:5000 LOGGING_LEVEL=debug KAFKA_WRITE_TOPIC_PARTITIONS=1 KAFKA_TOPIC_REPLICA_FACTOR=1 METRICS_ADDR=':8080' KAFKA_BROKERS=localhost:9092 KAFKA_READ_TOPICS=reclassified.fm.notification KAFKA_WRITE_TOPIC=active.alarm.change KAFKA_GROUP_ID=aalSbGroup DB_ADDRESSES=127.0.0.1 DB_SERVICE_USERNAME=aal DB_SERVICE_PASSWORD=aal REST_PORT=12806 TIME_FORMAT=RFC3339 DB_NAME=active_alarm_list CLEANUP_INTERVAL_IN_SEC=600 MAX_RETRY_COUNT=3 DB_REPLICATION_FACTOR="'datacenter1': 1" ./active-alarms-list

swagger-generate-server:
ifeq ($(SWAGGER), 1)
	swagger generate server -t . -f swagger.yml -A kmt
else
	echo "swagger version: $(CURR_SWAGGER_VER) is not available. Please install!"
endif