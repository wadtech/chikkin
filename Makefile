GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
SERVERBIN=chikkin

all: test build buildwin
build:
	$(GOBUILD) -o ./bin/$(SERVERBIN) -v ./cmd/gamey/main.go
buildwin:
	GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc $(GOBUILD) -o ./bin/$(SERVERBIN).exe -v ./cmd/gamey/main.go
test:
	$(GOTEST) -v ./...
testnocache:
	$(GOTEST) -v -count=1 ./...
clean:
	rm ./bin/$(SERVERBIN) ./bin/$(SERVERBIN).exe
deps:
	$(GOMOD) tidy
