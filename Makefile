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
	AIRTABLE_API_KEY="foo" \
	AIRTABLE_BASE_ID="foo" \
	MESSENGER_ACCESS_TOKEN="foo" \
	MESSENGER_VERIFY_TOKEN="foo" \
	MESSENGER_APP_SECRET="foo" \
	bin/airbot

.PHONY: all build bin clean run
