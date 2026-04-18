.PHONY: build run clean install

BINARY  := shepherd
CMD     := ./cmd/shepherd
LDFLAGS := -s -w

## build: Compile the binary to ./shepherd
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD)

## run: Build and run Shepherd
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
