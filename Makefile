TARGET := $(MAKECMDGOALS)
LIBDIR := $(CURDIR)/lib
export LIBDIR
DIRS := module protocol task
PWD  := $(shell pwd)

TMP := $(foreach n,$(DIRS),$(wildcard $(n)/*/Makefile))
SUBDIRS := $(subst /Makefile,,$(TMP))

all: $(SUBDIRS)
arm: $(SUBDIRS)
mac: $(SUBDIRS)
linux: $(SUBDIRS)
clean:$(SUBDIRS)

$(SUBDIRS):
	@echo $@
	$(MAKE) -C $@ M=$(PWD) $(TARGET)

.PHONY: all arm mac linux clean $(SUBDIRS)