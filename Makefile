GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
SERVERBIN=gamey

all: test build
build:
	$(GOBUILD) -o ./bin/$(SERVERBIN) -v ./cmd/gamey/main.go
test:
	$(GOTEST) -v ./...
testnocache:
	$(GOTEST) -v -count=1 ./...
clean:
	rm -f ./bin/$(SERVERBIN)
deps:
	$(GOMOD) tidy
