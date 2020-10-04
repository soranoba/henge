
build:
	go build

test:
	go test ./...
	cd tests; go test ./...

format:
	gofmt -w ./

lint:
	gofmt -d ./
	test -z $(shell gofmt -l ./)
