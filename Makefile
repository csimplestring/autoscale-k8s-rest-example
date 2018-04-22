
build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server -a -tags netgo

build-docker: build
	docker build . -t echo-api:latest
	rm server

tests:
	go clean -testcache
	go test ./...

check:
	go fmt ./...
	go vet ./...
	golint ./ ./customer

minikube:
	eval $(minikube docker-env)