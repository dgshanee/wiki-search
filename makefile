BINARY_NAME=wikisearch
default: build

fmt:
	go fmt ./...

vet:
	go vet ./...

build: fmt vet
	go build -o ${BINARY_NAME} .

run: build
	./${BINARY_NAME}

test: build
	go test ./...

clean:
	go clean
	rm -f ${BINARY_NAME}
