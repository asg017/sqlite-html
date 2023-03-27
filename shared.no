#include "sqlite3ext.h"

SQLITE_EXTENSION_INIT3

extern int go_sqlite3_extension_init(const char*, sqlite3*, char**);

int sqlite3_html_init(sqlite3* db, char** pzErrMsg, const sqlite3_api_routines *pApi) {
	SQLITE_EXTENSION_INIT2(pApi)
	return go_sqlite3_extension_init("html", db, pzErrMsg);
}