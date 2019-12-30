# These are the values we want to pass for VERSION and BUILD
# git tag 1.0.1
# git commit -am "One more change after the tags"
VERSION=1.0.0
BUILD=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'`
GIT_HASH=`git rev-parse HEAD`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-w -s -X github.com/jasonholmberg/slashspot/config.Version=${VERSION} -X github.com/jasonholmberg/slashspot/config.BuildTime=$(BUILD) -X github.com/jasonholmberg/slashspot/config.GitHash=$(GIT_HASH)"

test: 
	go test -cover ./... 

build: clean
	@echo Building Version: ${VERSION}
	@go build ${LDFLAGS} -o bin/slashspot-${VERSION} ./cmd/slashspot 

run: build 
	@echo "Make sure you set up the .env file"
	./bin/slashspot-${VERSION}

clean:
	rm -rf ./bin

.PHONY: all test build clean
