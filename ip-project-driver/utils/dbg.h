#ifndef DBG_H
#define DBG_H

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <stdarg.h>
#include <unistd.h>
#include <limits.h>

#define _NORMAL_    "\x1b[0m"
#define _BLACK_     "\x1b[30m"
#define _RED_       "\x1b[31m"
#define _GREEN_     "\x1b[32m"
#define _YELLOW_    "\x1b[33m"
#define _BLUE_      "\x1b[34m"
#define _MAGENTA_   "\x1b[35m"
#define _CYAN_      "\x1b[36m"
#define _WHITE_     "\x1b[37m"
#define _BRED_      "\x1b[1;31m"
#define _BGREEN_    "\x1b[1;32m"
#define _BYELLOW_   "\x1b[1;33m"
#define _BBLUE_     "\x1b[1;34m"
#define _BMAGENTA_  "\x1b[1;35m"
#define _BCYAN_     "\x1b[1;36m"
#define _BWHITE_    "\x1b[1;37m"

#define DBG_MODE(x)     (1ULL << (x))

#define DBG_ALL         (~0ULL)
#define DBG_ERROR       DBG_MODE(0)     /* error conditions             */
#define DBG_LINK        DBG_MODE(1)
#define DBG_NET         DBG_MODE(2)
#define DBG_ROUTE       DBG_MODE(3)

/// There can be no white space in these strings
#define DBG_NAMETAB \
{ "all", DBG_ALL }, \
{ "error", DBG_ERROR }, \
{ "link", DBG_LINK }, \
{ "net", DBG_NET }, \
{ "route", DBG_ROUTE }, \
{ NULL,         0 } 

#define DBG_COLORTAB \
{ DBG_ERROR, _RED_ }, \
{ DBG_LINK, _BYELLOW_ }, \
{ DBG_NET, _BBLUE_ }, \
{ DBG_ROUTE, _BMAGENTA_ }

#define dbg(mode, args...) do { \
	if (!dbg_initiated) dbg_init();  \
	if (dbg_modes & (mode)) { \
		char __buf__[LINE_MAX];\
		int __len__ = 0; \
		__len__ += sprintf(__buf__ + __len__, "%s", DCOLOR(mode));\
		__len__ += sprintf(__buf__ + __len__, ## args); \
		__len__ += sprintf(__buf__ + __len__, "%s", _NORMAL_);\
		__len__ = write(STDERR_FILENO, __buf__, __len__);\
		fflush(stderr);\
	} \
} while(0)

extern unsigned long long dbg_modes;
extern short dbg_initiated;
void dbg_init();
const char* DCOLOR(unsigned long long d_mode);

#endif
