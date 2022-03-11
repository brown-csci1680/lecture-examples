#include "dbg.h"

typedef struct dbg_mode {
  const char *d_name;
  unsigned long long d_mode;
} dbg_mode_t;

typedef struct dbg_mode_color {
  unsigned long long d_mode;
  const char *color;
} dbg_mode_color_t;

unsigned long long dbg_modes = 0;
short dbg_initiated = 0;

static dbg_mode_color_t dbg_colortab[] = {
  DBG_COLORTAB
};

static dbg_mode_t dbg_nametab[] = {
  DBG_NAMETAB
};

const char* DCOLOR(unsigned long long d_mode) {
    dbg_mode_color_t *mode;
    for (mode = dbg_colortab; mode->d_mode != 0; mode++) {
	if (mode->d_mode & d_mode)
	    return mode->color;
    }
    return _BWHITE_;
}

static void dbg_add_mode(const char *name) {
    int cancel;
    dbg_mode_t *mode;
    if (*name == '-') {
	cancel = 1;
	name++;
    } else cancel = 0;

    for (mode = dbg_nametab; mode->d_name != NULL; mode++)
	if (strcmp(name, mode->d_name) == 0) break;
    if (mode->d_name == NULL) {
	fprintf(stderr, "Warning: Unknown debug option: " "\"%s\"\n", name);
	return;
    }

    if (cancel) dbg_modes &= ~mode->d_mode;
    else dbg_modes = dbg_modes | mode->d_mode;
}

void dbg_init() {
    char env[256];
    char *name;
    const char *dbg_env;

    dbg_initiated = 1;
    dbg_modes = DBG_ERROR;
    dbg_env = getenv("DBG_MODES");
    if (!dbg_env) return;

    strncpy(env, dbg_env, sizeof(255));
    for (name = strtok(env,","); name; name = strtok(NULL, ","))
	dbg_add_mode(name);
}
