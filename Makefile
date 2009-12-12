CC=gcc
LD=gcc

all: virn

virn: virn.o
	$(LD) -o $@ -lapr-1 -laprutil-1 -lncurses $^ 

virn.o: virn.c
	$(CC) -g -c -I/usr/include/apr-1.0 $<

clean:
	rm *.o *~
