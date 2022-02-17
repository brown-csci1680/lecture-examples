#include <arpa/inet.h>
#include <inttypes.h>
#include <string.h>
#include <stdlib.h>
#include <stdio.h>

#include "dbg.h"
#include "parselinks.h"

lnxinfo_t *parse_links(char *filename) {
  FILE *f;
  lnxinfo_t *lnx;
  lnxbody_t *link, **next;
  int port, ret;
  char lhost[HOST_MAX_LENGTH], rhost[HOST_MAX_LENGTH], lv[INET_ADDRSTRLEN], rv[INET_ADDRSTRLEN];

	if ((f = fopen(filename, "r")) == NULL) {
		dbg(DBG_ERROR, "parse_links: fopen error\n");
		return NULL;
	}
	
	if (feof(f)) {
		fclose(f); 
		return NULL;
	}
  
	if ((lnx = malloc(sizeof(lnxinfo_t))) == NULL) {
		dbg(DBG_ERROR, "parse_links: malloc of lnxinfo_t error\n");
		fclose(f); 
		return NULL;
	}
	
	lnx->body = NULL;
	next = &lnx->body;
  if (fscanf(f, "%s %d", lhost, &port) != 2 || port < 0 || port > 0xffff) {
		dbg(DBG_ERROR, "parse_links: fscanf of local info error\n");
		free(lnx); 
		fclose(f); 
		return NULL;
	}
	lnx->local_phys_port = port;

	ret = 0;	
	while (!feof(f)) {
		ret = fscanf(f, "%s %d %s %s", rhost, &port, lv, rv);
		if (ret == EOF) {
			ret = 0;
			break;
		}

		link = malloc(sizeof(lnxbody_t));
    if (link == NULL) {
			dbg(DBG_ERROR, "parse_links: malloc of lnxbody_t error\n");
			ret = -1;
			break;
		}

		*next = link;
		link->next = NULL;
		next = &link->next;
	
		if (ret != 4 || port < 0 || port > 0xffff
				|| inet_pton(AF_INET, lv, &link->local_virt_ip) <= 0
				|| inet_pton(AF_INET, rv, &link->remote_virt_ip) <= 0) {
			dbg(DBG_ERROR, "parse_links: fscanf of link info error\n");
			ret = -1;
			break;
		}
		strcpy(link->remote_phys_host, rhost);
		link->remote_phys_port = port;
  }
	
	if (ret < 0) {
		fclose(f);
		free_links(lnx);
		return NULL;
	}
	return lnx;
}

void free_links(lnxinfo_t *lnx) {
  lnxbody_t *curr, *next;
	for (curr = lnx->body; curr != NULL; curr = next) {
		next = curr->next;
		free(curr);
	}
	free(lnx);
}
