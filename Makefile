.PHONY: all test test-server test-docker docker docker-clean publish-docker

REPO=github.com/alde/fusion
VERSION?=$(shell git describe --always HEAD | sed s/^v//)
DATE?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
DOCKERNAME?=alde/fusion
DOCKERTAG?=${DOCKERNAME}:${VERSION}
LDFLAGS=-X ${REPO}/version.Version=${VERSION} -X ${REPO}/version.BuildDate=${DATE}

SRC=$(shell find . -name '*.go')
TESTFLAGS="-v"

DOCKER_GO_SRC_PATH=/go/src/github.com/alde/fusion
DOCKER_GOLANG_RUN_CMD=docker run --rm -v "$(PWD)":$(DOCKER_GO_SRC_PATH) -w $(DOCKER_GO_SRC_PATH) golang:1.8 bash -c

PACKAGES=$(shell go list ./... | grep -v /vendor/)

all: test

glide:
	curl https://glide.sh/get | sh

deps:
	glide install

test: fusion
	go test ${TESTFLAGS} ${PACKAGES}

# Run tests cleanly in a docker container.
test-docker:
	$(DOCKER_GOLANG_RUN_CMD) "make glide test"

vet:
	go vet ${PACKAGES}

lint:
	go list ./... | grep -v /vendor/ | grep -v assets | xargs -L1 golint -set_exit_status

server/assets/assets.go: server/generate.go ${STATIC}
	go generate github.com/fusion-framework/fusion/server

fusion: ${SRC}
#	go build -ldflags "${LDFLAGS}" -o $@ github.com/alde/fusion/cmd/fusion
	go build -ldflags "${LDFLAGS}" -o $@ github.com/alde/fusion

docker/fusion: ${SRC}
#	CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -a -installsuffix cgo -o $@ github.com/alde/fusion/cmd/fusion
	CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -a -installsuffix cgo -o $@ github.com/alde/fusion

docker: docker/fusion docker/Dockerfile
	docker build -t ${DOCKERTAG} docker

docker-clean: docker/Dockerfile
	# Create the docker/fusion binary in the Docker container using the
	# golang docker image. This ensures a completely clean build.
	$(DOCKER_GOLANG_RUN_CMD) "make docker/fusion"
	docker build -t ${DOCKERTAG} docker

publish-docker:
#ifeq ($(strip $(shell docker images --format="{{.Repository}}:{{.Tag}}" $(DOCKERTAG))),)
#	$(warning Docker tag does not exist:)
#	$(warning ${DOCKERTAG})
#	$(warning )
#	$(error Cannot publish the docker image. Please run `make docker` or `make docker-clean` first.)
#endif
	docker push ${DOCKERTAG}
	git describe HEAD --exact 2>/dev/null && \
		docker tag ${DOCKERTAG} ${DOCKERNAME}:latest && \
		docker push ${DOCKERNAME}:latest || true
