UID:=$(shell id -u)
GO:=/usr/local/go/bin/go # default linux path

all: install

install:
ifneq ($(UID),0)
	@echo "Sorry, you are not root."
	@exit 1
endif
	@echo "Install 'nkb.service'"
	GOPATH=$(shell pwd)/../../../.. ${GO} build -v -o /usr/local/bin/nkb ./cmd/nkb
	cp ./nkb.service /etc/systemd/system/nkb.service
	systemctl daemon-reload
	systemctl enable nkb.service
	systemctl restart nkb.service
