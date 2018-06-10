COMMIT=$(shell git rev-parse --short HEAD)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X version.Commit=${COMMIT}"

all: build

bin:
	mkdir -p bin

build: bin
	go build ${LDFLAGS} -o bin/airbot ./cmd/airbot

clean:
	rm -rf bin

run: build
	GOOGLE_APPLICATION_CREDENTIALS=config/service-account.json \
	go run ${LDFLAGS} cmd/airbot/main.go

test:
	go test ./...

deploy:
	scripts/deploy

.PHONY: all build bin clean run test
