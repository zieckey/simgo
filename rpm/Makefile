SUBDIRS := 

OS=$(shell uname| tr a-z A-Z |cut -d_ -f1)
ifeq ($(OS), LINUX)
SUBDIRS += centos
endif # LINUX
ifeq ($(OS), FREEBSD)
SUBDIRS += freebsd
endif # FREEBSD

all:
	@for d in $(SUBDIRS); do \
        $(MAKE) -C $$d all; \
    done

clean:
	@for d in $(SUBDIRS); do \
		$(MAKE) -C $$d clean; \
	done

