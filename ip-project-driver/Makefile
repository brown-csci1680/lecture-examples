CC = gcc
CFLAGS = -g -Wall -Wextra -Werror -DREADLINE -Iutils -I.
LDFLAGS = -lpthread -lreadline

SRCS = node.c utils/dbg.c utils/parselinks.c utils/htable.c

all: example_node

example_node : $(SRCS)
	$(CC) $(CFLAGS) $^ -o $@ $(LDFLAGS)

clean:
	rm -f example_node
