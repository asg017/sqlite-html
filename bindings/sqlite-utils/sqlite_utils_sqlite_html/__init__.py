from sqlite_utils import hookimpl
import sqlite_html

from sqlite_utils_sqlite_html.version import __version_info__, __version__


@hookimpl
def prepare_connection(conn):
    conn.enable_load_extension(True)
    sqlite_html.load(conn)
    conn.enable_load_extension(False)
