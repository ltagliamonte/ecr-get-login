export GOBIN=$(shell pwd)
export GOPATH=$(shell pwd)/.go

all: build
build: $(GOBIN)/ecr-get-login

clean:
	rm -rf $(GOPATH) ecr-get-login

$(GOBIN)/ecr-get-login:
	mkdir -p $(GOPATH)
	go get -d .
	go build

test:
	test -z "$(shell gofmt -s -l main.go)"
	go vet main.go
