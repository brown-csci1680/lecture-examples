// server.c
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

#include "protocol.h"
#include "utils.h"

#define LISTEN_MAX 20


struct game_state {
    pthread_mutex_t lock;
    int32_t guess_number;
    int32_t total_guesses;

    // ...
};

struct game_state game_state;

// Data we want to track for each client
struct client_data {
    int sock;  // Socket file descriptor
    struct sockaddr_storage addr;
    socklen_t addr_size;

    pthread_t thread;
};

void init_game(void){
    game_state.guess_number = rand() % 1000;
    game_state.total_guesses = 0;

    printf("The number is %d\n", game_state.guess_number);
}

void *client_handler(void *data)
{
    struct client_data *cd = (struct client_data *)data;
    //printf("Client connected!\n");

    // Wait for a message from the client
    struct guess_message msg;
    recv_guess_message(cd->sock, &msg);
    // ^^ Not checking return value, this is bad!  What happens if there's an error?

    printf("Received guess:  %d\n", msg.number);

    int response = 0;

    // Modifications to server state start here
    pthread_mutex_lock(&game_state.lock);

    game_state.total_guesses++;
    if (msg.number < game_state.guess_number) {
	response = -1;
    } else if (msg.number > game_state.guess_number) {
	response = 1;
    } else {
	response = 0;
	// Client won, so restart

	init_game();
	// TODO:  This isn't great--the clients don't know the game has restarted!
	// What do we need to change about the server (and clients!) to change this?
    }

    // Done modifying server state
    pthread_mutex_unlock(&game_state.lock);

    send_guess_message(cd->sock, MESSAGE_TYPE_RESPONSE, response);


    pthread_exit(NULL);
}


int main(int argc, char **argv)
{
    // Goals for the server:
    //  - Listen for incoming connections, on some (address, port)
    //  - TCP connection, listen  on port 8888
    int rv;
    int s_listen;

    if (argc != 2) {
	printf("Usage:  server <port>\n");
	exit(1);
    }

    // Parse server port from command line
    char *server_port = argv[1];

    init_game();
    pthread_mutex_init(&game_state.lock, NULL);

    // Ask for a TCP socket that listens on all addresses
    struct addrinfo hints, *res;
    memset(&hints, 0, sizeof (struct addrinfo));
    hints.ai_family = AF_UNSPEC;     // IPv4 or IPv6 socket
    hints.ai_socktype = SOCK_STREAM; // TCP
    hints.ai_flags = AI_PASSIVE;     // Bind to all addresses on the system

    if ((rv = getaddrinfo(NULL, server_port, &hints, &res)) < 0) {
	perror("getaddrinfo");
	exit(1);
    }

    // Create the socket (s_listen)
    if((s_listen = socket(res->ai_family, res->ai_socktype, res->ai_protocol)) < 0) {
	perror("socket");
	exit(1);
    }

    // Bind my server to a specific address and port
    if((rv = bind(s_listen, res->ai_addr, res->ai_addrlen)) < 0) {
	perror("bind");
	exit(1);
    }

    // Listen for new connections
    if((rv = listen(s_listen, LISTEN_MAX)) < 0) {
	perror("listen");
	exit(1);
    }

    while(1) {
	int s_client;
	struct sockaddr_storage client_addr;
	socklen_t addr_size;

	// accept() blocks until a client has connected, and then
	// returns a file descriptor for a new socket that is used to
	// communicate directly with that client.  When this occurs,
	// information about the client's IP address is written to
	// client_addr and addr_size
	if ((s_client = accept(s_listen,
			       (struct sockaddr*)&client_addr, &addr_size)) < 0) {
	    perror("accept");
	}

	// Create the client data
	struct client_data *cd = (struct client_data *)malloc(sizeof(struct client_data));

	// Zero out and fill in client data structure
	memset(cd, 0, sizeof(struct client_data));
	cd->sock = s_client;
	memcpy(&cd->addr, &client_addr, sizeof(struct sockaddr_storage));
	cd->addr_size = addr_size;

	// Create a thread to handle the client request
	pthread_create(&cd->thread, NULL, client_handler, (void*)cd);
    }

    return 0;
}
