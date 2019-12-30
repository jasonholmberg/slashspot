.PHONY: all test clean

test: 
	go test -cover ./... 

build: clean
	go build -o bin/slashspot ./cmd/slashspot 

run: build 
	cp config/.env.template .env
	./bin/slashspot

clean:
	rm -rf ./bin
