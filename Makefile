# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_FOLDER=bin
BINARY_NAME=main
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_FILE=cmd/main.go

all: test build
build:
	$(GOBUILD) -o $(BINARY_FOLDER)/$(BINARY_NAME) -v $(MAIN_FILE)
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_FOLDER)/$(BINARY_NAME)
	rm -f $(BINARY_FOLDER)/$(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_FOLDER)/$(BINARY_NAME) -v $(MAIN_FILE)
	./$(BINARY_FOLDER)/$(BINARY_NAME)
deploy: clean test build-linux
	rsync -av . root@bento.do:/root/web && ssh 'root@bento.do' 'service bentoweb restart'

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_FOLDER)/$(BINARY_UNIX) -v $(MAIN_FILE)
