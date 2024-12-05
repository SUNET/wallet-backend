.PHONY : docker-build docker-push release

NAME 					:= vc
LDFLAGS                 := -ldflags "-w -s --extldflags '-static'"
LDFLAGS_DYNAMIC			:= -ldflags "-w -s"
CURRENT_BRANCH 			:= $(shell git rev-parse --abbrev-ref HEAD)

test: test-verifier

test-verifier:
	$(info Testing verifier)
	go test -v ./cmd/verifier

gosec:
	$(info Run gosec)
	gosec -color -nosec -tests ./...

staticcheck:
	$(info Run staticcheck)
	staticcheck ./...

start:
	$(info Run!)
	docker compose -f docker-compose.yaml up -d --remove-orphans

stop:
	$(info stopping VC)
	docker compose -f docker-compose.yaml rm -s -f

restart: stop start

get_release-tag:
	@date +'%Y%m%d%H%M%S%9N'

ifndef VERSION
VERSION := latest
endif


DOCKER_TAG_APIGW 		:= docker.sunet.se/dc4eu/apigw:$(VERSION)
DOCKER_TAG_VERIFIER		:= docker.sunet.se/dc4eu/verifier:$(VERSION)
DOCKER_TAG_REGISTRY 	:= docker.sunet.se/dc4eu/registry:$(VERSION)
DOCKER_TAG_PERSISTENT 	:= docker.sunet.se/dc4eu/persistent:$(VERSION)
DOCKER_TAG_GOBUILD 		:= docker.sunet.se/dc4eu/gobuild:$(VERSION)
DOCKER_TAG_MOCKAS 		:= docker.sunet.se/dc4eu/mockas:$(VERSION)
DOCKER_TAG_ISSUER 		:= docker.sunet.se/dc4eu/issuer:$(VERSION)
DOCKER_TAG_UI 			:= docker.sunet.se/dc4eu/ui:$(VERSION)
DOCKER_TAG_PORTAL 		:= docker.sunet.se/dc4eu/portal:$(VERSION)


build: proto build-verifier build-registry build-persistent build-mockas build-apigw build-ui

build-verifier:
	$(info Building verifier)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(NAME)_verifier ${LDFLAGS} ./cmd/verifier/main.go

build-registry:
	$(info Building registry)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(NAME)_registry ${LDFLAGS_DYNAMIC} ./cmd/registry/main.go

build-persistent:
	$(info Building persistent)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(NAME)_persistent ${LDFLAGS} ./cmd/persistent/main.go

build-mockas:
	$(info Building mockas)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(NAME)_mockas ${LDFLAGS} ./cmd/mockas/main.go

build-apigw:
	$(info Building apigw)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(NAME)_apigw ${LDFLAGS} ./cmd/apigw/main.go

build-ui:
	$(info Building ui)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(NAME)_ui ${LDFLAGS} ./cmd/ui/main.go

docker-build: docker-build-verifier docker-build-registry docker-build-persistent docker-build-mockas docker-build-apigw docker-build-issuer docker-build-ui docker-build-portal

docker-build-goland-debug: docker-build-verifier docker-build-registry docker-build-persistent docker-build-mockas docker-build-apigw docker-build-issuer docker-build-ui-goland-debug

docker-build-gobuild:
	$(info Docker Building gobuild with tag: $(VERSION))
	docker build --tag $(DOCKER_TAG_GOBUILD) --file dockerfiles/gobuild .

docker-build-verifier:
	$(info Docker Building verifier with tag: $(VERSION))
	docker build --build-arg SERVICE_NAME=verifier --tag $(DOCKER_TAG_VERIFIER) --file dockerfiles/worker .

docker-build-registry:
	$(info Docker Building registry with tag: $(VERSION))
	docker build --build-arg SERVICE_NAME=registry --tag $(DOCKER_TAG_REGISTRY) --file dockerfiles/worker .

docker-build-persistent:
	$(info Docker Building persistent with tag: $(VERSION))
	docker build --build-arg SERVICE_NAME=persistent --tag $(DOCKER_TAG_PERSISTENT) --file dockerfiles/worker .

docker-build-mockas:
	$(info Docker Building mockas with tag: $(VERSION))
	docker build --build-arg SERVICE_NAME=mockas --tag $(DOCKER_TAG_MOCKAS) --file dockerfiles/worker .

docker-build-apigw:
	$(info Docker building apigw with tag: $(VERSION))
	docker build --build-arg SERVICE_NAME=apigw --build-arg VERSION=$(VERSION) --tag $(DOCKER_TAG_APIGW) --file dockerfiles/worker .

docker-build-issuer:
	$(info Docker building issuer with tag: $(VERSION))
	docker build --build-arg SERVICE_NAME=issuer --tag $(DOCKER_TAG_ISSUER) --file dockerfiles/worker .

docker-build-ui:
	$(info Docker building ui with tag: $(VERSION))
	docker build --build-arg SERVICE_NAME=ui --tag $(DOCKER_TAG_UI) --file dockerfiles/ui_worker .

docker-build-ui-goland-debug:
	$(info Docker building ui with tag: $(VERSION))
	docker build --build-arg SERVICE_NAME=ui --tag $(DOCKER_TAG_UI) --file dockerfiles/ui_worker_goland_debug .

docker-build-portal:
	$(info Docker building portal with tag: $(VERSION))
	docker build --build-arg SERVICE_NAME=portal --tag $(DOCKER_TAG_PORTAL) --file dockerfiles/portal_worker .

docker-push-gobuild:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_GOBUILD)

docker-push-verifier:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_VERIFIER)

docker-push-registry:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_REGISTRY)

docker-push-mockas:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_MOCKAS)

docker-push-persistent:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_PERSISTENT)

docker-push-apigw:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_APIGW)

docker-push-issuer:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_ISSUER)

docker-push-ui:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_UI)

docker-push-portal:
	$(info Pushing docker images)
	docker push $(DOCKER_TAG_PORTAL)

docker-push: docker-push-verifier docker-push-registry docker-push-persistent docker-push-apigw docker-push-issuer docker-push-ui docker-push-mockas docker-push-portal
	$(info Pushing docker images)

docker-tag-apigw:
	$(info Tagging docker images)
	docker tag $(DOCKER_TAG_APIGW) docker.sunet.se/dc4eu/apigw:$(NEWTAG)

docker-tag-issuer:
	$(info Tagging docker images)
	docker tag $(DOCKER_TAG_ISSUER) docker.sunet.se/dc4eu/issuer:$(NEWTAG)

docker-tag-verifier:
	$(info Tagging docker images)
	docker tag $(DOCKER_TAG_VERIFIER) docker.sunet.se/dc4eu/verifier:$(NEWTAG)

docker-tag-registry:
	$(info Tagging docker images)
	docker tag $(DOCKER_TAG_REGISTRY) docker.sunet.se/dc4eu/registry:$(NEWTAG)

docker-tag-persistent:
	$(info Tagging docker images)
	docker tag $(DOCKER_TAG_PERSISTENT) docker.sunet.se/dc4eu/persistent:$(NEWTAG)

docker-tag-mockas:
	$(info Tagging docker images)
	docker tag $(DOCKER_TAG_MOCKAS) docker.sunet.se/dc4eu/mockas:$(NEWTAG)

docker-tag-ui:
	$(info Tagging docker images)
	docker tag $(DOCKER_TAG_UI) docker.sunet.se/dc4eu/ui:$(NEWTAG)

docker-tag: docker-tag-apigw docker-tag-issuer docker-tag-verifier docker-tag-registry docker-tag-persistent docker-tag-mockas docker-tag-ui
	$(info Tagging docker images)

check_current_branch:
	$(info Current branch: $(CURRENT_BRANCH))
ifeq ($(CURRENT_BRANCH),main)
	$(info main branch)
else
	$(error Not on main branch)
endif


release: check_current_branch
	$(info Release version: $(VERSION))
	git tag $(VERSION)
	git push origin ${VERSION}
	make docker-build
	make docker-push
	$(info Release version $(VERSION) done)
	$(info tag $(NEWTAG) from $(VERSION))
	make docker-tag
	make VERSION=$(NEWTAG) docker-push

	$(info point latest to $(NEWTAG))
	make NEWTAG=latest docker-tag
	make VERSION=latest docker-push

docker-pull:
	$(info Pulling docker images)
	docker pull $(DOCKER_TAG_APIGW)
	docker pull $(DOCKER_TAG_GOBUILD)
	docker pull $(DOCKER_TAG_MOCKAS)
	docker pull $(DOCKER_TAG_PERSISTENT)
	docker pull $(DOCKER_TAG_VERIFIER)
	docker pull $(DOCKER_TAG_REGISTRY)
	docker pull $(DOCKER_TAG_UI)

docker-archive:
	docker save --output docker_archives/vc_$(VERSION).tar $(DOCKER_TAG_VERIFIER) $(DOCKER_TAG_REGISTRY)


clean_redis:
	$(info Cleaning redis volume)
	docker volume rm vc_redis_data 

clean_docker_images:
	$(info Cleaning docker images)
	docker rmi $(DOCKER_TAG_VERIFIER) -f
	docker rmi $(DOCKER_TAG_REGISTRY) -f


ci_build: docker-build docker-push
	$(info CI Build)

proto: proto-status proto-registry proto-issuer

proto-registry:
	protoc --proto_path=./proto/ --go-grpc_opt=module=vc --go-grpc_out=. --go_opt=module=vc --go_out=. ./proto/v1-registry.proto

proto-status:
	protoc --proto_path=./proto/ --go-grpc_opt=module=vc --go_opt=module=vc --go_out=. --go-grpc_out=. ./proto/v1-status-model.proto 

proto-issuer:
	protoc --proto_path=./proto/ --go-grpc_opt=module=vc --go_opt=module=vc --go_out=. --go-grpc_out=. ./proto/v1-issuer.proto 

swagger: swagger-registry swagger-verifier swagger-apigw swagger-issuer swagger-fmt

swagger-fmt:
	swag fmt

swagger-registry:
	swag init -d internal/registry/apiv1/ -g client.go --output docs/registry --parseDependency --packageName docs

swagger-verifier:
	swag init -d internal/verifier/apiv1/ -g client.go --output docs/verifier --parseDependency --packageName docs

swagger-apigw:
	swag init -d internal/apigw/apiv1/ -g client.go --output docs/apigw --parseDependency --packageName docs

swagger-issuer:
	swag init -d internal/issuer/apiv1/ -g client.go --output docs/issuer --parseDependency --packageName docs

diagram:
	plantuml docs/diagrams/*.puml

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
		netcat-openbsd \
		plantuml
	$(info Install go packages)
	go install github.com/swaggo/swag/cmd/swag@latest && \
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
	go install golang.org/x/tools/cmd/deadcode@latest && \
	go install github.com/securego/gosec/v2/cmd/gosec@latest && \
	go install honnef.co/go/tools/cmd/staticcheck@latest