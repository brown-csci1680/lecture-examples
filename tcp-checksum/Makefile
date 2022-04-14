CC=gcc
CFLAGS=-Wall -O0 -g

.PHONY: all clean

all: tcpsum_example

tcpsum_example: tcpsum_example.c
	$(CC) $(CFLAGS) $^ -o $@

clean:
	rm -fv tcpsum_example
