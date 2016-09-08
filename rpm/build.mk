TEMP1 = $(shell echo `pwd`)
TEMP2 = $(shell echo `dirname ${TEMP1}`)
TEMP3 = $(shell echo `dirname ${TEMP2}`)
PROJ  = $(shell echo `basename ${TEMP3}`)

SRC			= $(PROJ).tar.gz
PKG 		= $(wildcard *.spec)
TEMP_DIR 	= /home/$(shell whoami)/tmp
ROOT_DIR 	= $(shell PWD=$$(pwd); echo $${PWD%%/$(PROJ)*}/$(PROJ))
SVNVERSION  = 1.0.$(shell svnversion ../../src)

