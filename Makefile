UID:=$(shell id -u)
GO:=/usr/local/go/bin/go # default linux path

all: install

install:
ifneq ($(UID),0)
	@echo "Sorry, you are not root."
	@exit 1
endif
	@echo "Install 'nkb.service'"
	${GO} install ./cmd/nkb
	@cp ./nkb.service /etc/systemd/system/nkb.service
	systemctl enable nkb
	sudo service nkb restart
