CC=gcc
CFLAGS=-O0 -Wall -g
LDFLAGS=

.PHONY: all clean

BINS=sender listener
all: $(BINS)

%.o: %.c
	$(CC) $(CFLAGS) $^ -c

sender: sender.o
	$(CC) $(LDLAGS) $^ -o $@

listener: listener.o utils.o
	$(CC) $(LDLAGS) $^ -o $@

clean:
	rm -rf *.o $(BINS)
