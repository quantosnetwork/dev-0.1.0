BINARY_NAME=quantosd
BINARY_PATH=bin
PROTO_PATH=proto
GO_PATH=${GOPATH}
export GODEBUG=gocacheverify=1

MAKE_WITH_MULTICPU = make -j${(`nproc`+1)}

fast:
	${MAKE_WITH_MULTICPU} build

all macos:
	make dep && make test && make build

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

diagram:
	go-plantuml -recursive . > quantos.puml

clean:
	rm -rf *.puml
	rm -rf *.png && rm -rf *.jpg
	rm -rf ${BINARY_PATH}/${BINARY_NAME}
	rm ./.version
	rm -rf vendor
	go clean -cache