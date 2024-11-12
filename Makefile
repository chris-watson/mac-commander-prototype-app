.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix 

.PHONY: test
test:
	go test ./...s

.PHONY: package
package:
	sh install/create_distro.sh

