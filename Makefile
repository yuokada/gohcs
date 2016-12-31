GOPATH=$(CURDIR)/.go
DEBUG=1

goget:
	@# echo $(GOPATH)
	@GOPATH=$(GOPATH) go get github.com/pkg/errors

build: goget
	cd src && GOPATH=$(GOPATH) go build -o gohealthcheck
