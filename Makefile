default: build

build:
	@cd jig && go get
	@cd jig && go build

install:
	@cd jig && go install

test:
	@go test
	@cd jig && go test

clean:
	@go clean
	@cd jig && go clean
	@rm -rf pkgroot

stage: build
	@mkdir -p pkgroot/usr/bin
	@cp jig/jig pkgroot/usr/bin/jig

deb: stage
	@fpm -C pkgroot -s dir -v $$(cat VERSION) -n jig -t deb usr

rpm: stage
	@fpm -C pkgroot -s dir -v $$(cat VERSION) -n jig -t rpm usr

