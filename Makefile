COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

all: build

bin:
	mkdir -p bin

build: bin
	go build ${LDFLAGS} -o bin/airbot .

clean:
	rm -rf bin

run: build
	ENV=development \
	PROJECT_ID=rising-artifact-182801 \
	BASE_ID=Shows \
	bin/airbot

deploy:
	scripts/deploy

.PHONY: all build bin clean run
