build: test
	go build -o troute main.go

test:
	go test ./...