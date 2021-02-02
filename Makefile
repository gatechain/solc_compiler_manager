#!/usr/bin/make -f
export GO111MODULE = on
export CGO_CFLAGS += -I$(SODIUM_PATH)/include
export CGO_LDFLAGS += -L$(SODIUM_PATH)/lib

# process build tags
whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

all: install lint

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

build-contract-tests-hooks:
ifeq ($(OS),Windows_NT)
	go build -mod=readonly $(BUILD_FLAGS) -o build/contract_tests.exe ./cmd/contract_tests
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/contract_tests ./cmd/contract_tests
endif

install: go.sum sodium
	go install $(BUILD_FLAGS) ./cmd/solc-compiler/

########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

clean:
	rm -rf build/

distclean: clean
	rm -rf vendor/

########################################
### Testing

test: test-unit test-build
test-all: test-race test-cover

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -timeout 0 -tags='ledger test_ledger_mock' ./...

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -timeout 0 -race -tags='ledger test_ledger_mock' ./...

test-cover:
	@go test -mod=readonly -timeout 0 -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

test-build: build
	@go test -mod=readonly -timeout 0 -p 4 `go list ./cli_test/...` -tags=cli_test -v

lint: golangci-lint
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

benchmark:
	@go test -mod=readonly -bench=. ./...

.PHONY: all build-linux install install-debug \
	go-mod-cache draw-deps clean build \
	setup-transactions setup-contract-tests-data start-gate run-lcd-contract-tests contract-tests \
	test test-all test-build test-cover test-unit test-race sodium
