CFLAGS=-I$(LUAJIT_SRC) -L$(LUAJIT_SRC) -lluajit

all: lib

lib:
	gcc -O2 -Wall -c lua_cjson.c strbuf.c fpconv.c  $(CFLAGS) && ar -cr libcjson.a *.o

clean:
	rm -rf *.a *.o

# export LUAJIT_SRC=luajit/src && make