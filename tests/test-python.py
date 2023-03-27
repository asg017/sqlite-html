import unittest
import sqlite3
import sqlite_html

class TestSqliteHtmlPython(unittest.TestCase):
  def test_path(self):
    db = sqlite3.connect(':memory:')
    db.enable_load_extension(True)

    self.assertEqual(type(sqlite_html.loadable_path()), str)
    
    sqlite_html.load(db)
    version, = db.execute('select html_version()').fetchone()
    self.assertEqual(version[0], "v")

if __name__ == '__main__':
    unittest.main()