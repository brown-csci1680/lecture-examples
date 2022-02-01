#ifndef __PROTOCOL_H__
#define __PROTOCOL_H__

#include <stdint.h>  // Provides uint8_t, int8_t, etc.

struct guess_message {
    uint8_t type;
    int32_t number;
} __attribute__((packed)); // Tell compiler to not add padding to this struct

#define MESSAGE_TYPE_GUESS 0
#define MESSAGE_TYPE_RESPONSE 1



int send_guess_message(int sock, uint8_t type, int32_t number);
int recv_guess_message(int sock, struct guess_message *msg);

#endif
