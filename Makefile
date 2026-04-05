.PHONY: generate build run clean

generate:
	go run github.com/a-h/templ/cmd/templ@latest generate

build: generate
	go build -o bin/server cmd/server/main.go

build-example: generate
	go build -o bin/example cmd/example/main.go

generate-example: build-example
	./bin/example

run: generate
	go run cmd/server/main.go

clean:
	rm -rf bin/
	rm -f drugprofile.db
	rm -f internal/views/*_templ.go
