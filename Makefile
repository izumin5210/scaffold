NAME := scaffold
VERSION := 0.0.1
REVISION := $(shell git describe --always)
LDFLAGS := -ldflags="-s -w -X \"main.Name=$(NAME)\" -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""
PACKAGE_DIRS := $(shell go list ./... 2> /dev/null | grep -v /vendor | grep -v /mock)
MOCK_FILES := $(shell git ls-files | grep _mock.go | paste -s -d "," -)

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
	@go get github.com/golang/mock/mockgen
	@go get github.com/golang/lint/golint
	@go get github.com/mitchellh/gox
	@go get github.com/mattn/goveralls
	@go get github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: generate
generate:
	go generate $(PACKAGE_DIRS)

.PHONY: build
build: generate
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: test-ci
test-ci: generate lint
	goveralls -service=travis-ci -ignore="$(MOCK_FILES)"

.PHONY: test
test: generate lint
	go test -v -timeout=30s -parallel=4 $(PACKAGE_DIRS)

.PHONY: lint
lint:
	go vet $(PACKAGE_DIRS)
	echo $(PACKAGE_DIRS) | xargs -n 1 golint -set_exit_status

.PHONY: package
package: clean deps
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

