CONTAINER_NAME=superman-detector
BINARY_NAME=superman-detector
GOOS=darwin
GOARCH=amd64

.PHONY: build
.PHONY: fmt
.PHONY: vet
.PHONY: docker.build
.PHONY: test.unit
.PHONY: test
.PHONY: clean

build: fmt vet superman-detector

fmt:
	go $@ ./...

vet:
	go $@ ./...

superman-detector:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-s -w" -o $(BINARY_NAME)

build.docker:
	docker build -t $(CONTAINER_NAME) .

test: test.unit

test.unit: fmt vet
	go test -v -race -cover ./...

clean:
	-rm $(BINARY_NAME)