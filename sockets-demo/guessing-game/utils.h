#ifndef __UTILS_H__
#define __UTILS_H__

#include <netdb.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>

void fatal_error(const char *message);

void *get_in_addr(struct sockaddr *sa);

void get_addr_string(struct sockaddr *sa, char *s, int len);
void print_addr(char *msg, struct sockaddr *sa);

#endif
