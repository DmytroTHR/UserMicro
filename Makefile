go-gen:
	@go generate -v ./main.go

go-build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/userservice