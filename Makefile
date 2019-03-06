TARGET := $(MAKECMDGOALS)
LIBDIR := $(CURDIR)/lib
export LIBDIR
DIRS := module protocol task
PWD  := $(shell pwd)

TMP := $(foreach n,$(DIRS),$(wildcard $(n)/*/Makefile))
SUBDIRS := $(subst /Makefile,,$(TMP))

$(shell mkdir -p ./lib/protocol)
$(shell mkdir -p ./lib/task)
$(shell mkdir -p ./lib/module)

all: $(SUBDIRS)
	go build main.go
arm: $(SUBDIRS)
	CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm  GOARM=7 go build main.go
clean:$(SUBDIRS)
	-rm main

$(SUBDIRS):
	@echo $@
	$(MAKE) -C $@ M=$(PWD) $(TARGET)

.PHONY: all arm mac linux clean $(SUBDIRS)