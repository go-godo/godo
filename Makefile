.PHONY: build get

build:
	@cd cmd/godo && go install

get:
	@go get -u github.com/go-godo/godo
	@go get -u github.com/go-godo/godo/cmd/godo

