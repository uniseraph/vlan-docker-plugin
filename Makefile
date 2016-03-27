SHELL = /bin/bash

TARGET       = vlan-docker-plugin
PROJECT_NAME = github.com/uniseraph/vlan-docker-plugin

MAJOR_VERSION = $(shell cat VERSION)
GIT_VERSION   = $(shell git log -1 --pretty=format:%h)
GIT_NOTES     = $(shell git log -1 --oneline)

BUILD_IMAGE     = uniseraph/vlan-docker-plugin:onbuild

IMAGE_NAME = uniseraph/vlan-docker-plugin
#REGISTRY = acs-reg.alipay.com


build-local:
	$(shell which godep) go build -a -v -ldflags "-B 0x$(shell head -c20 /dev/urandom|od -An -tx1|tr -d ' \n') -X ${PROJECT_NAME}/pkg/logging.ProjectName ${PROJECT_NAME} -X ${PROJECT_NAME}/cli.Version ${MAJOR_VERSION}(${GIT_VERSION})" -o ${TARGET}
	mkdir -p bundles/${MAJOR_VERSION}/binary
	mv ${TARGET} bundles/${MAJOR_VERSION}/binary
	@cd bundles/${MAJOR_VERSION}/binary && $(shell which md5sum) -b ${TARGET} | cut -d' ' -f1  > ${TARGET}.md5


build:
	docker build --rm -t ${BUILD_IMAGE} contrib/builder/binary
	docker run --rm -v $(shell pwd):/go/src/${PROJECT_NAME} -w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} make build-local

image:
	cp -r contrib/builder/image IMAGEBUILD && cp bundles/${MAJOR_VERSION}/binary/network-plugins IMAGEBUILD
	docker build --rm -t ${IMAGE_NAME}:${MAJOR_VERSION} IMAGEBUILD
	docker tag -f ${IMAGE_NAME}:${MAJOR_VERSION} ${IMAGE_NAME}:latest
	rm -rf IMAGEBUILD

push:
	docker tag -f ${IMAGE_NAME}:${MAJOR_VERSION} ${REGISTRY}/${IMAGE_NAME}:${MAJOR_VERSION}
	docker tag -f ${IMAGE_NAME}:latest ${REGISTRY}/${IMAGE_NAME}:latest
	docker push ${REGISTRY}/${IMAGE_NAME}:${MAJOR_VERSION}
	docker push ${REGISTRY}/${IMAGE_NAME}:latest



.PHONY: build build-local image
