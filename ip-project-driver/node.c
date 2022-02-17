#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <limits.h>
#include <inttypes.h>
#include <arpa/inet.h>
#include <errno.h>
#include <pthread.h>
#include <unistd.h>
#include <fcntl.h>

#ifdef READLINE
#include <readline/readline.h>
#include <readline/history.h>
#endif

#include "utils/dbg.h"
#include "utils/parselinks.h"

void help_cmd(const char *line) {
    (void) line;

    printf("- help, h: Print this list of commands\n"
	   "- interfaces, li: Print information about each interface, one per line\n"
	   "- routes, lr: Print information about the route to each known destination, one per line\n"
	   "- up [integer]: Bring an interface \"up\" (it must be an existing interface, probably one you brought down)\n"
	   "- down [integer]: Bring an interface \"down\"\n"
	   "- send [ip] [protocol] [payload]: sends payload with protocol=protocol to virtual-ip ip\n"
	   "- q: quit this node\n");
}

void interfaces_cmd(const char *line){
    (void) line;
    //TODO interfaces
    dbg(DBG_ERROR, "interfaces_cmd: NOT YET IMPLEMENTED\n");
}

void routes_cmd(const char *line){
    (void) line;
    //TODO routes
    dbg(DBG_ERROR, "routes_cmd: NOT YET IMPLEMENTED\n");
}

void down_cmd(const char *line){
    unsigned interface;

    if (sscanf(line, "down %u", &interface) != 1) {
	dbg(DBG_ERROR, "syntax error (usage: down [interface])\n");
	return;
    }
    //TODO down
    dbg(DBG_ERROR, "down_cmd: NOT YET IMPLEMENTED\n");
}

void up_cmd(const char *line){
    unsigned interface;

    if (sscanf(line, "up %u", &interface) != 1) {
	dbg(DBG_ERROR, "syntax error (usage: up [interface])\n");
	return;
    }
    //TODO up
    dbg(DBG_ERROR, "up_cmd: NOT YET IMPLEMENTED\n");
}

void send_cmd(const char *line){
    char ip_string[INET_ADDRSTRLEN];
    struct in_addr ip_addr;
    uint8_t protocol;
    int num_consumed;
    char *data;

    if (sscanf(line, "send %s %" SCNu8 "%n", ip_string, &protocol, &num_consumed) != 2) {
	dbg(DBG_ERROR, "syntax error (usage: send [ip] [protocol] [payload])\n");
	return;
    }

    if (inet_pton(AF_INET, ip_string, &ip_addr) == 0) {
	dbg(DBG_ERROR, "syntax error (malformed ip address)\n");
	return;
    }

    data = ((char *)line) + num_consumed + 1;

    if (strlen(data) < 1) {
	dbg(DBG_ERROR, "syntax error (payload unspecified)\n");
	return;
    }
    //TODO send
    dbg(DBG_ERROR, "send_cmd: NOT YET IMPLEMENTED\n");
}

struct {
  const char *command;
  void (*handler)(const char *);
} cmd_table[] = {
  {"help", help_cmd},
  {"h", help_cmd},
  {"interfaces", interfaces_cmd},
  {"i", interfaces_cmd},
  {"routes", routes_cmd},
  {"r", routes_cmd},
  {"down", down_cmd},
  {"up", up_cmd},
  {"send", send_cmd}
};


int main(int argc, char **argv){
    // (No arguments required for this example.  In the full version,
    // you should take an lnx file as input.)
    if (argc != 1) {
	dbg(DBG_ERROR, "usage: %s\n", argv[0]);
	return -1;
    }

#ifdef READLINE
    char* line;
    rl_bind_key('\t', rl_complete);
#else
    char line[LINE_MAX];
#endif
    char cmd[LINE_MAX];
    unsigned i;
    int ret;


    //TODO Initialize your layers!



    while (1) {
#ifdef READLINE
	if (!(line = readline("> "))) break;
#else
	dbg(DBG_ERROR, "> "); (void)fflush(stdout);
	if (!fgets(line, sizeof(line), stdin)) break;
	if (strlen(line) > 0 && line[strlen(line)-1] == '\n')
	    line[strlen(line)-1] = 0;
#endif

	ret = sscanf(line, "%s", cmd);
	if (ret != 1) {
	    if (ret != EOF) help_cmd(line);
	    continue;
	}
	if (!strcmp(cmd, "q")) break;

	for (i=0; i < sizeof(cmd_table) / sizeof(cmd_table[0]); i++){
	    if (!strcmp(cmd, cmd_table[i].command)){
		cmd_table[i].handler(line);
		break;
	    }
	}

	if (i == sizeof(cmd_table) / sizeof(cmd_table[0])){
	    dbg(DBG_ERROR, "error: no valid command specified\n");
	    help_cmd(line);
	    continue;
	}

#ifdef READLINE
	add_history(line);
	free(line);
#endif
    }


    //TODO Clean up your layers!


    printf("\nGoodbye!\n\n");
    return 0;
}
