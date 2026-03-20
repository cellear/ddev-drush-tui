BINARY=ddev-drush-tui

build:
	mkdir -p bin
	go build -o bin/$(BINARY) ./cmd/ddev-drush-tui

install: build
	mkdir -p $(HOME)/go/bin
	cp bin/$(BINARY) $(HOME)/go/bin/$(BINARY)

dist:
	mkdir -p dist
	GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY)-darwin-arm64 ./cmd/ddev-drush-tui
	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY)-darwin-amd64 ./cmd/ddev-drush-tui
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY)-linux-amd64 ./cmd/ddev-drush-tui

clean:
	rm -rf bin/ dist/
