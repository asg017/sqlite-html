import os
import sqlite3

from sqlite_html.version import __version_info__, __version__ 

def loadable_path():
  loadable_path = os.path.join(os.path.dirname(__file__), "html0")
  return os.path.normpath(loadable_path)

def load(conn: sqlite3.Connection)  -> None:
  conn.load_extension(loadable_path())
