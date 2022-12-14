DOC_EXECUTABLES = swag

test:
	CGO_ENABLED=0 go test -v ./...

coverage:
	go get golang.org/x/tools/cmd/cover
	go test -v -coverprofile cover.out ./... -mod=mod
	go tool cover -html=cover.out -o cover.html

build:
	go mod vendor
	GO111MODULE=on
	GOFLAGS="-mod=vendor"
	CGO_ENABLED=0 go build -o faceit-backend-test ./cmd/main.go

run: build
	./faceit-backend-test

check-doc-reqs:
	$(foreach bin,$(DOC_EXECUTABLES),\
		$(if $(shell command -v $(bin) 2> /dev/null),$(info Found `$(bin)`),$(error Please install `$(bin)`)))

generate-doc: check-doc-reqs
	echo "Generating swagger files"
	cd cmd;swag init -o ../docs --parseDependency --parseInternal