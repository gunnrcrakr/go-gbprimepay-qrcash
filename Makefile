build:
	go build -o bin/main main.go

run: build
	./bin/main

watch:
	reflex -s -r '\.go$$' make run