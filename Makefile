.PHONY: build

build:
	@go install
	@cd cmd/gosu && go install

