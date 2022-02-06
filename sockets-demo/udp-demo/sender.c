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

//#include "utils.h"


#define ADDR_STRING_SIZE 64
#define INPUT_SIZE 1024


int main(int argc, char **argv)
{
    int rv;
    int sock;

    if (argc != 4) {
	printf("Listener:  sender <address> <port> <message>\n");
	exit(1);
    }

    // Parse server port from command line
    char *dest_address = argv[1];
    char *dest_port = argv[2];
    char *message = argv[3];

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

    if ((rv = getaddrinfo(dest_address, dest_port, &hints, &servinfo)) != 0) {
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

	break;
    }

    if (res == NULL) {
	printf("Could not bind to socket\n");
	exit(1);
    }


    struct sockaddr *dest_sa = res->ai_addr;
    socklen_t dest_addr_len = res->ai_addrlen;

    int bytes_sent;

    if ((bytes_sent = sendto(sock, message, strlen(message), 0,
			     dest_sa, dest_addr_len)) < 0) {
	perror("sendto");
	exit(1);
    }
    printf("Sent %d bytes\n", bytes_sent);

    close(sock);
    freeaddrinfo(servinfo);

    return 0;
}
