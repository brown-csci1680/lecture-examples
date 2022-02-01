#include <errno.h>
#include <netdb.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <sys/types.h>

void fatal_error(const char *message)
{
    perror(message);
    exit(1);
}

void *get_in_addr(struct sockaddr *sa)
{
    if (sa->sa_family == AF_INET) {
        return &(((struct sockaddr_in*)sa)->sin_addr);
    }

    return &(((struct sockaddr_in6*)sa)->sin6_addr);
}

#define ADDR_MAX 256
void print_addr(char *msg, struct sockaddr *sa)
{
    char addr_str[ADDR_MAX];
    switch(sa->sa_family) {
    case AF_INET:
	inet_ntop(AF_INET, &(((struct sockaddr_in *)sa)->sin_addr),
		  addr_str, ADDR_MAX);
	break;

    case AF_INET6:
	inet_ntop(AF_INET6, &(((struct sockaddr_in6 *)sa)->sin6_addr),
		  addr_str, ADDR_MAX);
	break;

    default:
	strncpy(addr_str, "Unknown AF", ADDR_MAX);
	break;
    }

    printf("%s%s\n", (msg != NULL) ? msg : "", addr_str);
}
