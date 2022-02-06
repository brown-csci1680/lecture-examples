CC=gcc
CFLAGS=-O0 -Wall -g

.PHONY: all clean

BINS=server client
LDFLAGS=
all: $(BINS)

%.o: %.c
	$(CC) $(CFLAGS) $^ -c

client: client.o utils.o protocol.o
	$(CC) $(LDFLAGS) $^ -o $@

server: server.o utils.o protocol.o
	$(CC) $(LDFLAGS) $^ -lpthread -o $@

clean:
	rm -rf *.o $(BINS)
