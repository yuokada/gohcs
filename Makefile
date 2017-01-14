GOPATH=$(CURDIR)/.go
DEBUG=1

goget:
	@# echo $(GOPATH)
	@GOPATH=$(GOPATH) go get github.com/pkg/errors

build: goget
	cd src && GOPATH=$(GOPATH) go build -o gohealthcheck

install: build
	# UNAME := $(shell uname)
	# ifeq ( $(UNAME), Linux)
	# 	echo "install process"
	# endif
	@echo "install phase"
	install -o root -g root -m 0755 src/gohealthcheck /usr/local/bin/gohealthcheck
	install -o root -g root -m 0644 etc/gohealthcheck.conf /etc/tmpfiles.d/gohcheck.conf
  install -o root -g root -m 0644 etc/gohealthcheck.service /etc/systemd/system/gohealthcheck.service
