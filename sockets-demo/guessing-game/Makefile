CC=gcc
CFLAGS=-O0 -Wall -g

.PHONY: all clean

BINS=server client
all: $(BINS)


client: client.c utils.c protocol.c
	$(CC) $(CFLAGS) $^ -o $@

server: server.c utils.c protocol.c
	$(CC) $(CFLAGS) $^ -lpthread -o $@

clean:
	rm -rf *.o $(BINS)
