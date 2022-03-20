CC=gcc
CFLAGS=-Wall -O3 -g

.PHONY: all clean

all: ipsum_example

ipsum_example: ipsum_example.c
	$(CC) $(CFLAGS) $^ -o $@

clean:
	rm -fv ipsum_example
