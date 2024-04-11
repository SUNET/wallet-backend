NAME 					:= vc
LDFLAGS                 := -ldflags "-w -s --extldflags '-static'"
LDFLAGS_DYNAMIC			:= -ldflags "-w -s"

test: test-verifier test-datastore

test-verifier:
	$(info Testing verifier)
	go test -v ./cmd/verifier

test-datastore:
	$(info Testing datastore)
	go test -v ./cmd/datastore

gosec:
	$(info Run gosec)
	gosec -color -nosec -tests ./...

staticcheck:
	$(info Run staticcheck)
	staticcheck ./...

start:
	$(info Run!)
	docker-compose -f docker-compose.yaml up -d --remove-orphans

stop:
	$(info stopping VC)
	docker-compose -f docker-compose.yaml rm -s -f

restart: stop start

get_release-tag:
	@date +'%Y%m%d%H%M%S%9N'

ifndef VERSION
VERSION := latest
endif

docker-build: docker-build-verifier docker-build-datastore docker-build-registry docker-build-cache docker-build-persistent

DOCKER_TAG_APIGW 		:= docker.sunet.se/dc4eu/apigw:$(VERSION)
DOCKER_TAG_VERIFIER		:= docker.sunet.se/dc4eu/verifier:$(VERSION)
DOCKER_TAG_DATASTORE	:= docker.sunet.se/dc4eu/datastore:$(VERSION)
DOCKER_TAG_REGISTRY 	:= docker.sunet.se/dc4eu/registry:$(VERSION)
DOCKER_TAG_CACHE 		:= docker.sunet.se/dc4eu/cache:$(VERSION)
DOCKER_TAG_PERSISTENT 	:= docker.sunet.se/dc4eu/persistent:$(VERSION)
DOCKER_TAG_GOBUILD 		:= docker.sunet.se/dc4eu/gobuild:$(VERSION)
DOCKER_TAG_MOCKAS 		:= docker.sunet.se/dc4eu/mockas:$(VERSION)


build: proto build-verifier build-datastore build-registry build-cache build-persistent build-mockas build-apigw


build-verifier:
	$(info Building verifier)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(NAME)_verifier ${LDFLAGS} ./cmd/verifier/main.go

build-datastore:
	$(info Building datastore)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(NAME)_datastore ${LDFLAGS} ./cmd/datastore/main.go

build-registry:
	$(info Building registry)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(NAME)_registry ${LDFLAGS_DYNAMIC} ./cmd/registry/main.go

build-cache:
	$(info Building cache)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(NAME)_cache ${LDFLAGS} ./cmd/cache/main.go

build-persistent:
	$(info Building persistent)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(NAME)_persistent ${LDFLAGS} ./cmd/persistent/main.go

build-mockas:
	$(info Building mockas)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(NAME)_mockas ${LDFLAGS} ./cmd/mockas/main.go

build-apigw:
	$(info Building apigw)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(NAME)_apigw ${LDFLAGS} ./cmd/apigw/main.go


docker-build: docker-build-verifier docker-build-datastore docker-build-registry docker-build-cache docker-build-persistent docker-build-mockas docker-build-apigw

docker-build-gobuild:
	$(info Docker Building gobuild with tag: $(VERSION))
	docker build --tag $(DOCKER_TAG_GOBUILD) --file dockerfiles/gobuild .

docker-build-verifier:
	$(info Docker Building verifier with tag: $(VERSION))
	docker build --build-arg SERVICE_NAME=verifier --tag $(DOCKER_TAG_VERIFIER) --file dockerfiles/worker .

docker-build-datastore:
	$(info Docker Building datastore with tag: $(VERSION))
	docker build --build-arg SERVICE_NAME=datastore --tag $(DOCKER_TAG_DATASTORE) --file dockerfiles/worker .

docker-build-registry:
	$(info Docker Building registry with tag: $(VERSION))
	docker build --build-arg SERVICE_NAME=registry --tag $(DOCKER_TAG_REGISTRY) --file dockerfiles/worker .

docker-build-cache:
	$(info Docker Building cache with tag: $(VERSION))
	docker build --build-arg SERVICE_NAME=cache --tag $(DOCKER_TAG_CACHE) --file dockerfiles/worker .

docker-build-persistent:
	$(info Docker Building persistent with tag: $(VERSION))
	docker build --build-arg SERVICE_NAME=persistent --tag $(DOCKER_TAG_PERSISTENT) --file dockerfiles/worker .

docker-build-mockas:
	$(info Docker Building mockas with tag: $(VERSION))
	docker build --build-arg SERVICE_NAME=mockas --tag $(DOCKER_TAG_MOCKAS) --file dockerfiles/worker .

docker-build-apigw:
	$(info Docker building apigw with tag: $(VERSION))
	docker build --build-arg SERVICE_NAME=apigw --tag $(DOCKER_TAG_APIGW) --file dockerfiles/worker .

docker-push-gobuild:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_GOBUILD)

docker-push-datastore:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_DATASTORE)

docker-push-verifier:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_VERIFIER)

docker-push-registry:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_REGISTRY)

docker-push-mockas:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_MOCKAS)

docker-push-cache:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_CACHE)

docker-push-persistent:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_PERSISTENT)

docker-push-apigw:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_APIGW)

docker-push: docker-push-datastore docker-push-datastore docker-push-verifier docker-push-registry docker-push-cache docker-push-persistent docker-push-apigw
	$(info Pushing docker images)

docker-pull:
	$(info Pulling docker images)
	docker pull $(DOCKER_TAG_APIGW)
	docker pull $(DOCKER_TAG_GOBUILD)
	docker pull $(DOCKER_TAG_MOCKAS)
	docker pull $(DOCKER_TAG_CACHE)
	docker pull $(DOCKER_TAG_PERSISTENT)
	docker pull $(DOCKER_TAG_VERIFIER)
	docker pull $(DOCKER_TAG_DATASTORE)
	docker pull $(DOCKER_TAG_REGISTRY)

docker-archive:
	docker save --output docker_archives/vc_$(VERSION).tar $(DOCKER_TAG_VERIFIER) $(DOCKER_TAG_DATASTORE) $(DOCKER_TAG_REGISTRY)


clean_redis:
	$(info Cleaning redis volume)
	docker volume rm vc_redis_data 

clean_docker_images:
	$(info Cleaning docker images)
	docker rmi $(DOCKER_TAG_VERIFIER) -f
	docker rmi $(DOCKER_TAG_DATASTORE) -f
	docker rmi $(DOCKER_TAG_REGISTRY) -f


ci_build: docker-build docker-push
	$(info CI Build)

release-tag:
	git tag -s ${RELEASE} -m"release ${RELEASE}"

release_push:
	git push --tags

release: release-tag release_push
	$(info making release ${RELEASE})

proto: proto-status proto-registry

proto-registry:
	protoc --proto_path=./proto/ --go-grpc_opt=module=vc --go-grpc_out=. --go_opt=module=vc --go_out=. ./proto/v1-status-model.proto ./proto/v1-registry.proto

proto-status:
	protoc --proto_path=./proto/ --go-grpc_opt=module=vc --go_opt=module=vc --go_out=. --go-grpc_out=. ./proto/v1-status-model.proto 

swagger: swagger-registry swagger-datastore swagger-verifier swagger-apigw swagger-fmt

swagger-fmt:
	swag fmt

swagger-registry:
	swag init -d internal/registry/apiv1/ -g client.go --output docs/registry --parseDependency --packageName docs

swagger-datastore:
	swag init --exclude ./vendor/ -d internal/datastore/apiv1/ -g client.go --output docs/datastore --parseDependency --packageName docs

swagger-verifier:
	swag init -d internal/verifier/apiv1/ -g client.go --output docs/verifier --parseDependency --packageName docs

swagger-apigw:
	swag init -d internal/apigw/apiv1/ -g client.go --output docs/apigw --parseDependency --packageName docs

openapi3-apigw: swagger-apigw
	docker run --rm -v ./docs/apigw:/work openapitools/openapi-generator-cli:latest-release generate -i /work/swagger.yaml -o ./work/v3 -g openapi-yaml --minimal-update
	docker run --rm -v ./docs/apigw:/work openapitools/openapi-generator-cli:latest-release generate -i /work/swagger.json -o ./work/v3 -g openapi --minimal-update
	sudo mv ./docs/apigw/v3/openapi/openapi.yaml ./internal/apigw/httpserver/openapiv3/openapi.yaml
	sudo mv ./docs/apigw/v3/openapi.json ./internal/apigw/httpserver/openapiv3/openapi.json
	sudo rm -rf ./docs/apigw/v3
	sudo chown -R $(shell id -u):$(shell id -g) ./internal/apigw/httpserver/openapiv3/

install-tools:
	$(info Install from apt)
	apt-get update && apt-get install -y \
		protobuf-compiler \
		netcat-openbsd
	$(info Install from go)
	go install github.com/swaggo/swag/cmd/swag@latest && \
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


clean-apt-cache:
	$(info Cleaning apt cache)
	rm -rf /var/lib/apt/lists/*

vscode:
	$(info Install APT packages)
	sudo apt-get update && sudo apt-get install -y \
		protobuf-compiler \
		netcat-openbsd
	$(info Install go packages)
	go install github.com/swaggo/swag/cmd/swag@latest && \
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
	go install golang.org/x/tools/cmd/deadcode@latest && \
	go install github.com/securego/gosec/v2/cmd/gosec@latest && \
	go install honnef.co/go/tools/cmd/staticcheck@latest && \
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest