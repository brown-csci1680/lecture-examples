#include <stdio.h>
#include <pthread.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>

#include "htable.h"

static htable_node_t *__htable_lookup( htable_t *ht, unsigned int id );

void htable_init( htable_t *ht, unsigned int cap ) {
	unsigned int i;

	ht->ht_hash = (list_t*) malloc( sizeof( list_t ) * cap );
	ht->ht_cap = cap;
	ht->ht_size = 0;

	for( i = 0; i < cap; i++ )
		list_init( &ht->ht_hash[ i ] );
}

void htable_destroy( htable_t *ht ) {
	unsigned int i;
	htable_node_t *hn;

	for( i = 0; i < ht->ht_cap; i++ ) {
		list_iterate_begin( &ht->ht_hash[ i ], hn, htable_node_t, hn_link ) {
			free( hn );
		} list_iterate_end();
	}

	free( ht->ht_hash );
}

void *htable_get( htable_t *ht, unsigned int id ) {
	htable_node_t *hn;

	if( ( hn = __htable_lookup( ht, id ) ) ) return hn->hn_data;
	else return NULL;
}

void *htable_put( htable_t *ht, unsigned int id, void *data ) { 
	htable_node_t *hn;
	void *old = NULL;

	if( !( hn = __htable_lookup( ht, id ) ) ) {
		hn = (htable_node_t*) malloc( sizeof( htable_node_t ) );
		hn->hn_id = id;
		list_insert_head( &ht->ht_hash[ id % ht->ht_cap ], &hn->hn_link );
		ht->ht_size++;
	} else old = hn->hn_data;

	hn->hn_data = data;

	return old;
}

void *htable_remove( htable_t *ht, unsigned int id ) {
	htable_node_t *hn;

	if( ( hn = __htable_lookup( ht, id ) ) ) {
		void *data = hn->hn_data;
		list_remove( &hn->hn_link );
		free( hn );
		ht->ht_size--;
		return data;
	} else return NULL;
}

htable_node_t *__htable_lookup( htable_t *ht, unsigned int id ) {
	htable_node_t *hn;
	
	list_iterate_begin( &ht->ht_hash[ id % ht->ht_cap ], hn, htable_node_t, hn_link ) {
		if( hn->hn_id == id ) return hn;
	} list_iterate_end();

	return NULL;
}
