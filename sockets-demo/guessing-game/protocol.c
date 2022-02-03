// protocol.c

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


// Create and send a guess over a socket
int send_guess_message(int sock, uint8_t type, int32_t number)
{
    int b;

#if 0
    // Method 1:  call send() for each datatype

    if ((b = send(sock, &type, sizeof(uint8_t), 0)) <= 0) {
	perror("send");
	return -1;
    }

    // When we send multi-byte data, we use network byte order, which is BIG ENDIAN
    // Functions to convert byte order:  ntohs, ntohl, htons, htonl
    int32_t number_to_send = htonl(number);
    if ((b = send(sock, &number_to_send, sizeof(int32_t), 0)) <= 0) {
	perror("send");
	return -1;
    }
#endif

    // Method 2:  Build and pack a struct with the message
    struct guess_message msg;
    msg.type = type;

    // When we send multi-byte data, our protocol requires we use
    // network byte order, which is BIG ENDIAN
    // Functions to convert byte order:  ntohs, ntohl, htons, htonl
    msg.number = htonl(number);

    if ((b = send(sock, &msg, sizeof(struct guess_message), 0)) <= 0) {
	perror("send");
	return -1;
    }

    // There's more to do here!  What happens if send() sends fewer bytes than we requested?
    // What happens when the server closes the connection?

    return 0;
}

int recv_all(int sock, char *buffer, int total_size)
{
    int total_bytes_read = 0;
    int to_read = total_size;

    char *ptr = buffer;

    while (to_read > 0) {
	int bytes_read = recv(sock, ptr, to_read, 0);
	if (bytes_read <= 0) {
	    perror("recv");
	    return -1;
	}

	to_read -= bytes_read;
	ptr += bytes_read;
	total_bytes_read += bytes_read;
    }

    return total_bytes_read;
}

int recv_guess_message(int sock, struct guess_message *msg)
{
    // recv() may receive fewer bytes than we want!  Instead, we use a wrapper
    // to ensure we always read the appropriate number of bytes
    int b = recv_all(sock, (char *)msg, sizeof(struct guess_message));

    // Something went wrong, abort
    // What happens when the client just closes the connection vs. an error occurs?
    if (b <= 0) {
	return -1;
    }

    // Need to reformat message for host
    msg->number = ntohl(msg->number);

    // error handling might happen here...
    //printf("Received %d bytes\n", b);

    return 1;
}
