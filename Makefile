.PHONY: all test build clean

test: 
	go test -cover ./... 

build: clean
	go build -o bin/slashspot ./cmd/slashspot 

run: build 
	echo "Make sure you set up the .env file"
	./bin/slashspot

clean:
	rm -rf ./bin
