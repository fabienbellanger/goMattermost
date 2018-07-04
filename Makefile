# Go parameters
GOCMD=go
GOINSTALL=$(GOCMD) install
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=goSendNotification
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

install:
	$(GOINSTALL) ./...

runMattermost:
	$(GORUN) main.go mattermost -r project/ -p .
mattermost: install runMattermost

runMattermostNoDB:
	$(GORUN) main.go mattermost -r project/ -p . --no-database
mattermostNoDB: install runMattermostNoDB

build: 
	$(GOBUILD) -o $(BINARY_NAME) -v

test: 
	$(GOTEST) -cover ./...

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run-prod:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	$(GOGET) -u github.com/spf13/cobra/cobra
	# $(GOGET) -u github.com/labstack/echo/...
	$(GOGET) -u github.com/go-sql-driver/mysql
	$(GOGET) -u github.com/fatih/color
