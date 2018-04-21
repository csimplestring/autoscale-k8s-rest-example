
build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server -a -tags netgo

build-docker: build
	docker build . -t echo-rest/api

tests:
	go clean -testcache
	go test ./...
