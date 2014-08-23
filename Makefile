.PHONY: build get

build:
	@go install
	@cd cmd/gosu && go install
	@cd util && go install

get:
	@go get -u github.com/mgutz/gosu
	@go get -u github.com/mgutz/gosu/cmd/gos

