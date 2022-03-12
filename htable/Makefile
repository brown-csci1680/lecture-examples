CC=gcc
CFLAGS=-Wall -O3 -g

.PHONY: all clean

all: htable_example

htable_example: htable.c htable_example.c
	$(CC) $(CFLAGS) $^ -o $@

clean:
	rm fv htable_example
