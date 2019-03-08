CFLAGS=-I$(LUAJIT_SRC)  -L$(LUAJIT_SRC) -lluajit

all: lib

lib:
	gcc -O2 -Wall -c *.c  $(CFLAGS) && ar -cr liblpeg.a *.o

clean:
	rm -rf *.a *.o

# export LUAJIT_SRC=luajit/src && make
