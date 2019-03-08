CFLAGS=-I$(LUAJIT_SRC) -L$(LUAJIT_SRC) -lluajit

all: lib

lib:
	gcc -O2 -Wall -fPIC -c *.c  $(CFLAGS) && ar -cr libpbc.a *.o

clean:
	rm -rf *.a *.o

# export LUAJIT_SRC=luajit/src && make
