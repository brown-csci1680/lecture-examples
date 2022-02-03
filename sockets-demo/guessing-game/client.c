// client.c

#include <netdb.h>
#include <errno.h>
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <pthread.h>
#include <arpa/inet.h>
#include <sys/types.h>
#include <sys/socket.h>

#include "protocol.h"
#include "utils.h"

#define SERVER_PORT "9999"
#define MAX_READ_SIZE 1024

int main(int argc, char **argv)
{
    int sock, rv;
    struct addrinfo hints, *res, *servinfo;

    if (argc != 4) {
	printf("Usage:  client <address> <port> <guess>\n");
	exit(1);
    }

    char *server_address = argv[1];
    char *server_port = argv[2];
    char *guess_str = argv[3];
    int32_t guess_number = (int32_t)atoi(guess_str);

    memset(&hints, 0, sizeof(hints));
    hints.ai_family = AF_UNSPEC; // IPv4, IPv6 socket
    hints.ai_socktype = SOCK_STREAM; // TCP socket

    // Get an address struct for the server
    // In doing so, getaddrinfo() will resolve server_address
    // (eg. "1.2.3.4", "localhost", "cs.brown.edu", ...) into addrinfo structs
    // which can be used for other socket system calls.
    if ((rv = getaddrinfo(server_address, server_port,
			  &hints, &servinfo)) != 0) {
	perror(gai_strerror(rv));
	exit(1);
    }

    // For reasons we will discuss, there may be multiple possible resolutions for
    // different socket types, IP addresses, so here we iterate over these possibilities,
    // connect to each one, and pick the first one that works
    // For more info see https://beej.us/guide/bgnet/html/#getaddrinfoman
    for (res = servinfo; res != NULL; res = res->ai_next) {
	if ((sock = socket(res->ai_family, res->ai_socktype, res->ai_protocol)) < 0) {
	    //perror("socket"); // Not an error!  (Just need to try the next one)
	    continue;
	}

	// Try to connect to the specified address
	if (connect(sock, res->ai_addr, res->ai_addrlen) < 0) {
	    //perror("connect"); // Not an error! (Just need to keep trying)
	    close(sock);
	    continue;
	}

	break;
    }

    // If we exhausted all of the possibilities without getting a response, THIS is an error
    if (res == NULL) {
	printf("Failed to connect\n");
	exit(1);
    }

    // After this point, our socket has been created!


    // Otherwise, we can send our guess to the server
    send_guess_message(sock, MESSAGE_TYPE_GUESS, guess_number);

    // TODO next class:  Wait for a reply...
    struct guess_message response;
    rv = recv_guess_message(sock, &response);
    if (rv == -1) {
	printf("Error receiving response\n");
    } else {
	if (response.number > 0) {
	    printf("Wrong!  The number is higher than %d\n", guess_number);
	} else if (response.number < 0) {
	    printf("Wrong!  The number is lower than %d\n", guess_number);
	} else {
	    printf("Yay, you won!!\n");
	}
    }

    // Cleanup
    close(sock);
    freeaddrinfo(res);

    return 0;
}
