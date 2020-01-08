all:
	CGO_ENABLED=0 go build -o /usr/local/bin/ioa ./cmd/main.go

buildDocker:
	docker build -t cyjme/ioa -f ./Dockerfile .

