CC=gcc
CFLAGS=-O2 -Wall -g
LDFLAGS=-lpthread

.PHONY: all clean

BINS=server client
all: $(BINS)

%.o: *.c
	$(CC) $(CFLAGS) -c $^

client: client.o utils.o
	$(CC) $(LDFLAGS) $^ -o $@

server: server.o utils.o
	$(CC) $(LDFLAGS) $^ -o $@

clean:
	rm -rf *.o $(BINS)
