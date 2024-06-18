run:
	go run main.go files/example.lox

run_prompt:
	go run main.go

test:
	go generate
	go test ./...