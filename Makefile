.PHONY: build run clean install

BINARY  := baa
CMD     := ./cmd/baa
LDFLAGS := -s -w

## build: Compile the binary to ./baa
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD)

## run: Build and run BAA
run: build
	./$(BINARY)

## install: Install to $GOPATH/bin
install:
	go install -ldflags "$(LDFLAGS)" $(CMD)

## clean: Remove build artifacts
clean:
	rm -f $(BINARY)

## vet: Run static analysis
vet:
	go vet ./...

## help: Show this help
help:
	@grep -E '^## ' Makefile | sed 's/## //' | column -t -s ':'
