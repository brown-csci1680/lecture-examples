// server.c
#include <netdb.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <pthread.h>
#include <sys/types.h>
#include <sys/socket.h>

#define SERVER_PORT "9999"
#define LISTEN_MAX 20

struct client_data {
    int sock;
    struct sockaddr_storage addr;
    socklen_t addr_size;
    pthread_t thread;
};

void *client_handler(void *client_data)
{
    struct client_data *client = (struct client_data*)client_data;

    printf("Hello client!\n");

    close(client->sock);
    free(client);
    pthread_exit(NULL);
}

int main(int argc, char **argv)
{
    int rv;
    int s_listen;
    struct addrinfo hints, *res;

    memset(&hints, 0, sizeof(struct addrinfo));
    hints.ai_family = AF_UNSPEC;
    hints.ai_socktype = SOCK_STREAM;
    hints.ai_flags = AI_PASSIVE; // Pick any address

    if ((rv = getaddrinfo(NULL, SERVER_PORT, &hints, &res)) < 0) {
	perror(gai_strerror(rv));
    }

    if ((s_listen = socket(res->ai_family, res->ai_socktype, res->ai_protocol)) < 0) {
	perror("socket");
    }

    if ((rv = bind(s_listen, res->ai_addr, res->ai_addrlen)) < 0) {
	perror("bind");
    }

    if ((rv = listen(s_listen, LISTEN_MAX)) < 0) {
	perror("listen");
    }

    for(;;) {
	struct sockaddr_storage client_addr;
	socklen_t addr_size;
	int s_client;

	if ((s_client = accept(s_listen, (struct sockaddr *)&client_addr, &addr_size)) < 0) {
	    perror("accept");
	}

	// Create the client data
	struct client_data *client = (struct client_data*)malloc(sizeof(struct client_data));
	memset(client, 0, sizeof(struct client_data));
	memcpy(&client->addr, &client_addr, addr_size);
	client->addr_size = addr_size;

	pthread_create(&client->thread, NULL, client_handler, (void*)client);
    }

    return 0;
}
