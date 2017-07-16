NAME := scaffold
VERSION := 0.0.1
REVISION := $(shell git describe --always)
LDFLAGS := -ldflags="-s -w -X \"main.Name=$(NAME)\" -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

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
test-ci: lint vet test-race
	goveralls -service=travis-ci

.PHONY: test
test: lint vet test-race
	go test -v -timeout=30s -parallel=4 ./...

.PHONY: test-race
test-race:
	go test -race .

.PHONY: vet
vet:
	go vet *.go

.PHONY: lint
lint:
	@go get github.com/golang/lint/golint
	golint ./...

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
