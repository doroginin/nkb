UID:=$(shell id -u)
GO:=/usr/local/go/bin/go # default linux path

all: install

install:
ifneq ($(UID),0)
	@echo "Sorry, you are not root."
	@exit 1
endif
	@echo "Install 'keywatcher.service'"
	${GO} install
	@cp ./keywatcher.service /etc/systemd/system/keywatcher.service
	systemctl enable keywatcher
	sudo service keywatcher restart
