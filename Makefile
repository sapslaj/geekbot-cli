build: build/geekbot-cli
build/geekbot-cli: tidy download
	go build -v -o build/geekbot-cli .

.PHONY: clean download install tidy

clean:
	rm -f build/geekbot-cli

download:
	go mod download

install: build/geekbot-cli
	mkdir -p $$HOME/.local/bin
	cp build/geekbot-cli $$HOME/.local/bin

tidy:
	go mod tidy
