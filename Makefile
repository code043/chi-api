build:
	@go build -o bin/api main.go

test:
	@go test -v ./...
	
run: build
	@./bin/api