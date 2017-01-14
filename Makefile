GOPATH=$(CURDIR)/.go
DEBUG=1

goget:
	@# echo $(GOPATH)
	@GOPATH=$(GOPATH) go get github.com/pkg/errors

build: goget
	cd src && GOPATH=$(GOPATH) go build -o gohcs

install: build
	# UNAME := $(shell uname)
	# ifeq ( $(UNAME), Linux)
	# 	echo "install process"
	# endif
	@echo "install phase"
	# make directories
	install -o root -g root -m 0775 -d /var/run/gohcs
	install -o root -g root -m 0775 -d /etc/gohcs
	# copy files
	install -o root -g root -m 0755 src/gohcs /usr/local/bin/gohcs
	install -o root -g root -m 0644 etc/gohcs.conf /etc/tmpfiles.d/gohcheck.conf
	install -o root -g root -m 0644 etc/gohcs.service /etc/systemd/system/gohcs.service
