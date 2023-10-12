.PHONY: build go-build css-build run watch-css

# Build both Go and CSS
build: go-build css-build

# Build the Go project
go-build:
	go build -o tmp/main

# Build the CSS using TailwindCSS
css-build:
	tailwindcss -i web/input.css -o web/public/output.css --minify

# Run the Go server
run:
	go run main.go

# Watch the CSS for changes with TailwindCSS
watch-css:
	tailwindcss -i web/input.css -o web/public/output.css --watch

