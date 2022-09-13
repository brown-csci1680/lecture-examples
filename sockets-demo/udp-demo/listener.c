// listener.c
#include <netdb.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <unistd.h>
#include <pthread.h>
#include <sys/types.h>
#include <sys/socket.h>

#include "utils.h"


#define ADDR_STRING_SIZE 64
#define INPUT_SIZE 1024


int main(int argc, char **argv)
{
    int rv;
    int sock;

    if (argc != 2) {
	printf("Listener:  server <port>\n");
	exit(1);
    }

    // Parse server port from command line
    char *server_port = argv[1];

    // Ask for a socket that listens on all addresses
    struct addrinfo hints, *res, *servinfo;
    memset(&hints, 0, sizeof (struct addrinfo));
    hints.ai_family = AF_INET;       // Request an IPv4 socket
                                     // For simplicity, we standardize on IPv4 here
                                     // UDP has no connect() like with TCP, so we can't
                                     // follow the same procedure with the TCP client where
                                     // we tried both IPv4 and IPv6
    hints.ai_socktype = SOCK_DGRAM;  // UDP socket
    hints.ai_flags = AI_PASSIVE;     // Bind to all addresses on the system

    if ((rv = getaddrinfo(NULL, server_port, &hints, &servinfo)) != 0) {
	perror("getaddrinfo");
	exit(1);
    }

    // Look at all the results and bind to the first one
    // (Technically, we should be able to eliminate this loop, since we only picked AF_INET)
    for (res = servinfo; res != NULL; res = res->ai_next) {
	if((sock = socket(res->ai_family, res->ai_socktype, res->ai_protocol)) < 0) {
	    //perror("socket");
	    continue;
	}

	// Bind bind to the address and port
	if((rv = bind(sock, res->ai_addr, res->ai_addrlen)) < 0) {
	    close(sock);
	    continue;
	}

	break;
    }

    if (res == NULL) {
	printf("Could not bind to socket\n");
	exit(1);
    }


    while(1) {
	struct sockaddr_storage their_addr;
	socklen_t addr_size = sizeof(struct sockaddr_storage);
	int bytes_received;
	char buffer[INPUT_SIZE];

	if ((bytes_received = recvfrom(sock, buffer, INPUT_SIZE - 1, 0,
				       (struct sockaddr *)&their_addr, &addr_size)) == -1) {
	    perror("recvfrom");
	    exit(1);
	}

	// Null-terminate whatever we received for safety
	buffer[bytes_received] = '\0';

	char addr_string[ADDR_STRING_SIZE];
	get_addr_string((struct sockaddr *)&their_addr,
			addr_string, ADDR_STRING_SIZE);

	printf("Received %d bytes from %s:  %s\n",
	       bytes_received, addr_string, buffer);
    }

    close(sock);
    freeaddrinfo(servinfo);

    return 0;
}
