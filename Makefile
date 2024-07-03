run:
	go run main.go files/example.lox

run_prompt:
	go run main.go

gen_tokens:
	go run token/gen/main.go token

test:
	go generate
	go test ./...