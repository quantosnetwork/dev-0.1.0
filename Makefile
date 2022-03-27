BINARY_NAME=quantosd
BINARY_PATH=bin
PROTO_PATH=proto

build:
	go build -o ${BINARY_PATH}/${BINARY_NAME}

test:
	go test ./...
test_coverage:
	go test  ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all