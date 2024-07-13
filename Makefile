run:
	go run main.go files/example.lox

benchmark:
	go run main.go files/fib.lox

run_prompt:
	go run main.go

cleanup:
	rm token/token.go 2> /dev/null

gen_tokens: cleanup
	go run token/gen/main.go token

test:
	go generate
	go test ./...