APP_NAME ?= Commander
PLATFORM ?= darwin
ARCH ?= amd64

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix 

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	@echo "Building the application..."
	@GOOS=$(PLATFORM) GOARCH=$(ARCH) go build -o ./build/$(APP_NAME) ./cmd/main.go
	@test -f ./build/$(APP_NAME) || (echo "Build failed"; exit 1)

.PHONY: package
package: build
ifeq ($(PLATFORM),darwin)
	echo "Entering package target..."
	echo "PLATFORM is set to $(PLATFORM)"
	echo "Running packaging script..."
	sh install/create_$(PLATFORM)_$(ARCH)_distro.sh $(APP_NAME)
else
	$(error Packaging is currently only supported on darwin platform)
endif

.PHONY: run
run:
	go run ./cmd/main.go