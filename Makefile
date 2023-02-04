APP:=$(notdir $(patsubst %/,%, $(CURDIR)))

all: deps test build

.PHONY: deps
deps:
	go install github.com/nicksnyder/go-i18n/v2/goi18n@v2
	go mod download

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	CGO_ENABLED=0 go build -o $(APP)

.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(APP)

# Get all localized strings
.PHONY: lang-extract
lang-extract:
	goi18n extract

# Validate that all strings are translated
.PHONY: lang-validate
lang-validate:
	goi18n merge active.*.toml

# Merge new translations in the active files
.PHONY: lang-merge
lang-merge:
	goi18n merge active.*.toml translate.*.toml
