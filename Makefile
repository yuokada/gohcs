GOPATH=$(CURDIR)/.go
DEBUG=1
PREFIX:=""
USER=$(shell whoami)
GROUP=$(shell groups)
.PHONY: rpm

goget:
	@# echo $(GOPATH)
	@GOPATH=$(GOPATH) go get github.com/pkg/errors

build: goget
	cd src && GOPATH=$(GOPATH) go build -o gohcs

install:
	@echo "install phase"
	install -o ${USER} -g ${GROUP} -m 0775 -d ${PREFIX}/var/run/gohcs
	install -o ${USER} -g ${GROUP} -m 0775 -d ${PREFIX}/etc/gohcs
	install -o ${USER} -g ${GROUP} -m 0775 -d ${PREFIX}/etc/tmpfiles.d
	install -o ${USER} -g ${GROUP} -m 0775 -d ${PREFIX}/etc/systemd/system/
	install -o ${USER} -g ${GROUP} -m 0775 -d ${PREFIX}/usr/local/bin

	install -o ${USER} -g ${GROUP} -m 0755 src/gohcs         ${PREFIX}/usr/local/bin/gohcs

	install -o ${USER} -g ${GROUP} -m 0644 etc/tmpfiles.d/gohcs.conf    ${PREFIX}/etc/tmpfiles.d/gohcs.conf
	install -o ${USER} -g ${GROUP} -m 0644 etc/systemd/system/gohcs.service ${PREFIX}/etc/systemd/system/gohcs.service
	install -o ${USER} -g ${GROUP} -m 0644 etc/gohcs/checklist.json     ${PREFIX}/etc/gohcs/checklist.json

rpm:
	/bin/bash ./buildrpm.sh
#	 # if you want to check specfile, aadd --review-spec
#	       checkinstall \
#	 --fstrans=no \
#	 --install=no \
#	 -R \
#	 -A x86_64 \
#	 --pkglicense=MIT \
#	 --pakdir=pkg \
#	 --pkgversion=1.0.0 \
#	 --delspec=no \
#	 --review-spec \
#	 --backup=no \
#	  -y
