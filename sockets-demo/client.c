// client.c

#include <netdb.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <pthread.h>
#include <arpa/inet.h>
#include <sys/types.h>
#include <sys/socket.h>

#include "utils.h"

#define SERVER_PORT "9999"
#define MAX_READ_SIZE 1024

int main(int argc, char **argv)
{
    int sock, rv, bytes_read;
    struct addrinfo hints, *res;

    memset(&hints, 0, sizeof(hints));
    hints.ai_family = AF_UNSPEC;
    hints.ai_socktype = SOCK_STREAM;

    if ((rv = getaddrinfo("localhost", SERVER_PORT, &hints, &res)) != 0) {
	fatal_error(gai_strerror(rv));
    }

    if ((sock = socket(res->ai_family, res->ai_socktype, res->ai_protocol)) < 0) {
	fatal_error("socket");
    }

    char addr_str[256];
    inet_ntop(res->ai_family, get_in_addr((struct sockaddr *)res->ai_addr), addr_str, sizeof(addr_str));
    printf("Connecting to %s...\n", addr_str);

    if (connect(sock, res->ai_addr, res->ai_addrlen) < 0) {
	fatal_error("connect");
    }

    printf("Connected...\n");


    char buffer[MAX_READ_SIZE];
    if ((bytes_read = recv(sock, buffer, MAX_READ_SIZE - 1, 0)) == -1) {
	fatal_error("recv");
    }

    buffer[bytes_read] = '\0';
    printf("Received %d bytes:  %s\n", bytes_read, buffer);

    close(sock);
    freeaddrinfo(res);

    return 0;
}
