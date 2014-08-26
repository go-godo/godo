.PHONY: build get

build:
	@cd cmd/gosu && go install

get:
	@go get -u github.com/mgutz/gosu
	@go get -u github.com/mgutz/gosu/cmd/gosu

