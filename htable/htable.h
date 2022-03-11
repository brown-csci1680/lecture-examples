#ifndef __HASHTABLE_H__
#define __HASHTABLE_H__

#include <pthread.h>

#include "list.h"

/* FIXME make this a doubly-hashed, dynamically groweable hashtable */

typedef struct htable {
	list_t *ht_hash;				/* table entries */
	unsigned int ht_size;				    /* table size */
	unsigned int ht_cap;				        /* table capacity */
} htable_t;

typedef struct htable_node {
	list_t hn_link;			/* link */
	unsigned int hn_id;				        /* hash id */
	void *hn_data;					/* data */
} htable_node_t;

void htable_init( htable_t *ht, unsigned int cap );
void htable_destroy( htable_t *ht );
void *htable_get( htable_t *ht, unsigned int id );
void *htable_put( htable_t *ht, unsigned int id, void *data );
void *htable_remove( htable_t *ht, unsigned int id );

#define htable_iterate_begin( ht, key, var, type ) \
do { \
    unsigned int ___i; \
    htable_t *__ht = (ht); \
    htable_node_t *__hnode; \
	for(___i = 0;___i < __ht->ht_cap;___i++ ) { \
		list_iterate_begin( &__ht->ht_hash[___i ], __hnode, htable_node_t, hn_link ) { \
			(var) = (type*) __hnode->hn_data; \
			(key) = __hnode->hn_id; \
            do

#define htable_iterate_end() \
			while( 0 ); \
	    } list_iterate_end(); \
	} \
} while( 0 )

#endif /* __HASHTABLE_H__ */
