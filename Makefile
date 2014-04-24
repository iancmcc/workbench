default: build

build:
	@cd jig && go build

install:
	@cd jig && go install

test:
	@go test
	@cd jig && go test

clean:
	@go clean
	@cd jig && go clean
