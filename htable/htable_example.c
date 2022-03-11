#include <stdio.h>
#include <stdlib.h>
#include <assert.h>

#include <arpa/inet.h>

#include "htable.h"

// Max size for our hash table
// (This doesn't allocate memory, it just wraps the hash values)
#define MY_TABLE_MAX 1024

// Example data structure
struct thing {
    long x;
    long y;

    htable_node_t htable_node;
};

int main(int argc, char **argv)
{
    htable_t htable;

    // Initialize hash table
    htable_init(&htable, MY_TABLE_MAX);

    // Make some elements
    struct thing t1 = {42, 5};
    struct thing t2 = {1, 2};
    struct thing t3 = {1024, 2};

    // Add elements to the hash table
    // id is the hash value (type unsigned int)
    // Note:  If you want to store something other than an integer type,
    // you'll need to make a hash function.  However, note that:
    // 1. An IP address is just a 32-bit number (see below)
    // 2. If you needed a type larger than an "unsigned int",
    //   feel free to modify htable.{c,h} to use, say, an unsigned long
    htable_put(&htable, 0, &t1);
    htable_put(&htable, 1, &t2);
    htable_put(&htable, 2, &t3);

    // Get a specific element
    struct thing *result = (struct thing *)htable_get(&htable, 0);
    if (result != NULL) {
	printf("Got element (%ld, %ld)\n",
	       result->x, result->y);
    }

    // Replace an element
    struct thing t4 = {100, 200};

    assert(htable.ht_size == 3);
    htable_put(&htable, 0, &t4); // Replace element with ID 0
    assert(htable.ht_size == 3); // Size hasn't changed

    // Iterate over the table
    unsigned int key;
    struct thing *value;
    printf("---- Iterating over hash table ----\n");
    htable_iterate_begin(&htable, key, value, struct thing) {
	printf("ID:  %d, Thing:  {%ld, %ld}\n",
	       key, value->x, value->y);
    } htable_iterate_end();
    printf("---- Done iterating ----\n");

    // For other methods (htable_remove, htable_destroy), see htable.{c,h}

    // ***********************************************
    // ***** Want to use an IP address as a key? *****
    // ***********************************************

    // First, let's just get a sockaddr struct like you're
    // likely to have in the assignment
    int rv;
    struct sockaddr_in addr; // IPv4 sockaddr struct
                             // (be sure to create all of your sockets with AF_INET)
    if ((rv = inet_pton(AF_INET, "127.0.0.1", &addr)) < 0) {
	perror("inet_pton");
    }

    // Now, we can get the address in 32-bit form from the sockaddr struct
    // For details, see:  https://beej.us/guide/bgnet/html/#structsockaddrman
    uint32_t addr_as_int = addr.sin_addr.s_addr;

    // Make some element and store it into the hash table
    struct thing t = {5, 5};
    htable_put(&htable, addr_as_int, &t);

    // Now let's get it back out
    result = (struct thing *)htable_get(&htable, addr_as_int);
    if (result != NULL) {
	printf("Got element (%ld, %ld)\n",
	       result->x, result->y);
    }


    return 0;
}
