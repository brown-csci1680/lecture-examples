#ifndef _PARSELINKS_H_
#define _PARSELINKS_H_

#include <stddef.h>
#include <netinet/if_ether.h>
#include <netinet/in.h>

#define HOST_MAX_LENGTH 256 // RFC 2181

typedef struct lnxbody_s {
	char remote_phys_host[HOST_MAX_LENGTH]; 
	uint16_t remote_phys_port;
	struct in_addr local_virt_ip; // Network-byte-order
	struct in_addr remote_virt_ip; // Network-byte-order

	struct lnxbody_s *next;
} lnxbody_t;

typedef struct {
	uint16_t local_phys_port;
	lnxbody_t* body;
} lnxinfo_t;

// Returns a lnxinfo_t struct. Call free_links to clean up
lnxinfo_t *parse_links(char *filename);
void free_links(lnxinfo_t *lnx);

#endif
