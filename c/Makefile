DIR=$(shell pwd)

LUADIR=$(DIR)/LuaJIT-2.1.0-beta1

all: lib

lib:
	cd $(LUADIR) && make \
	&& export LUAJIT_SRC=$(LUADIR)/src \
	&& cd $(DIR)/pbc && make \
	&& cd $(DIR)/lpeg && make \
	&& cd $(DIR)/cjson && make

clean:
	cd $(DIR)/pbc && make clean \
	&& cd $(DIR)/lpeg && make clean \
	&& cd $(DIR)/cjson && make clean \
	&& cd $(LUADIR) && make clean
	