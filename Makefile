utest:
	export GOPATH=$(shell pwd)
	go test -v sdk/* 2>&1

deploy:
	export GOPATH=$(shell pwd)
	go build -v sdk/
build:
	export GOPATH=$(shell pwd)
	go build -o main *go

test:
	${MAKE} utest


.PHONY: test utest deploy
