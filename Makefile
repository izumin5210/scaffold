NAME := scaffold
VERSION := 0.0.1
REVISION := $(shell git describe --always)
LDFLAGS := -ldflags="-s -w -X \"main.Name=$(NAME)\" -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""
PACKAGE_DIRS   := $(shell go list ./... 2> /dev/null | grep -v /vendor/)

XC_ARCH := 386 amd64
XC_OS := darwin linux windows

.PHONY: clean
clean:
	rm -rf bin pkg

.PHONY: clobber
clobber: clean
	rm -rf vendor

.PHONY: deps
deps:
	@go get github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: build
build:
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: test-ci
test-ci: lint
	goveralls -service=travis-ci

.PHONY: test
test: lint
	go test -v -timeout=30s -parallel=4 $(PACKAGE_DIRS)

.PHONY: lint
lint:
	go vet $(PACKAGE_DIRS)
	@go get github.com/golang/lint/golint
	echo $(PACKAGE_DIRS) | xargs -n 1 golint -set_exit_status

.PHONY: package
package: clean deps
	@go get github.com/mitchellh/gox
	gox \
		$(LDFLAGS) \
		-parallel=5 \
		-os="$(XC_OS)" \
		-arch="$(XC_ARCH)" \
		-output="pkg/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)"
	@for pkg in pkg/*;do \
		tar zcf $${pkg}.tar.gz $${pkg}; \
		rm -rf $${pkg}; \
	done

.PHONY: release
release: package
	ghr --username $GITHUB_USER --token $GITHUB_TOKEN --replace --prerelease --debug pre-release pkg/
