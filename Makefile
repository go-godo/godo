.PHONY: build get

build:
	@cd cmd/godo && go install

get:
	@go get -u gopkg.in/godo.v1
	@go get -u gopkg.in/godo.v1/cmd/godo

