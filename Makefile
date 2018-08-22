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

runNotification:
	$(GORUN) main.go notification -r project/ -p .
notification: install runNotification

runNotificationNoDB:
	$(GORUN) main.go notification -r project/ -p . --no-database -a mattermost -b master
notificationNoDB: install runNotificationNoDB

runApi:
	$(GORUN) main.go web
serve: install runApi

runMail:
	$(GORUN) main.go mail
mail: install runMail

runDb:
	$(GORUN) main.go db
db: install runDb

runDbInit:
	$(GORUN) main.go db --init
dbInit: install runDbInit

runDbDump:
	$(GORUN) main.go db --dump
dbDump: install runDbDump

runConfig:
	$(GORUN) main.go config
config: install runConfig

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
	$(GOGET) -u github.com/labstack/echo/...
	$(GOGET) -u github.com/go-sql-driver/mysql
	$(GOGET) -u github.com/fatih/color
	$(GOGET) -u github.com/dgrijalva/jwt-go
