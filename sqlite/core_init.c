#include "sqlite-html.h"
int core_init(const char *dummy) {
  return sqlite3_auto_extension((void *) sqlite3_html_init);
}
