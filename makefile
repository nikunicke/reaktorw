all:
	go build -o reaktor-warehouse ./cmd/reaktorw

test:
	go test ./...