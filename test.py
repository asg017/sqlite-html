import sqlite3
import unittest
import os

EXT_PATH = "dist/html0"

db = sqlite3.connect(":memory:")

db.execute("create table fbefore as select name from pragma_function_list")
db.execute("create table mbefore as select name from pragma_module_list")

db.enable_load_extension(True)
db.load_extension(EXT_PATH)

db.execute("create temp table fafter as select name from pragma_function_list")
db.execute("create temp table mafter as select name from pragma_module_list")

class TestHtml(unittest.TestCase):
  def test_funcs(self):
    funcs = list(map(lambda a: a[0], db.execute("select name from fafter where name not in (select name from fbefore) order by name").fetchall()))
    self.assertEqual(funcs, [
      "html",
      "html_attr_get",
      "html_attr_has",
      "html_attribute_get",
      "html_attribute_has",
      "html_count",
      "html_debug",
      "html_element",
      "html_escape",
      "html_extract",
      "html_table",
      "html_text",
      "html_trim",
      "html_unescape",
      "html_version",
    ])
  
  def test_modules(self):
    funcs = list(map(lambda a: a[0], db.execute("select name from mafter where name not in (select name from mbefore) order by name").fetchall()))
    self.assertEqual(funcs, [
      "html_each",
    ])

  def test_http_version(self):
    v, = db.execute("select html_version()").fetchone()
    self.assertEqual(v, "v0.0.0")
  
  def test_html_debug(self):
    d, = db.execute("select html_debug()").fetchone()
    lines = d.splitlines()
    self.assertEqual(len(lines), 4)
    self.assertTrue(lines[0].startswith("Version"))
    self.assertTrue(lines[1].startswith("Commit"))
    self.assertTrue(lines[2].startswith("Runtime"))
    self.assertTrue(lines[3].startswith("Date"))

  def test_html_table(self):
    d, = db.execute("select html_table('a')").fetchone()
    self.assertEqual(d, "<table>a")
  
  def test_html_escape(self):
    d, = db.execute("select html_escape('<a>')").fetchone()
    self.assertEqual(d, "&lt;a&gt;")
  
  def test_html_unescape(self):
    d, = db.execute("select html_unescape('&lt;a')").fetchone()
    self.assertEqual(d, "<a")
  
  def test_html_trim(self):
    a,b = db.execute("""select html_trim('  a '), html_trim('
    bb 
    ')""").fetchone()
    self.assertEqual(a, "a")
    self.assertEqual(b, "bb")
  
  def test_html(self):
    a, b, c = db.execute("select html('a'), html('<p>b'), html('<ohno');").fetchone()
    self.assertEqual(a, "a")
    self.assertEqual(b, "<p>b</p>")
    # TODO what should this do
    self.assertEqual(c, None)
  
  def test_html_element(self):
    a, b, c = db.execute("""select 
      html_element('p', null, "ayoo"), 
      html_element('p', json_object("class", "cool"), "ayoo"), 
      html_element(
        'p',
        null,
        html_element("span", null, "My name is "),
        html_element("b", null, "Alex Garcia"),
        "."
      )
    """).fetchone()
    self.assertEqual(a, "<p>ayoo</p>")
    self.assertEqual(b, "<p class=\"cool\">ayoo</p>")
    # TODO what should this do
    self.assertEqual(c, "<p><span>My name is </span><b>Alex Garcia</b>.</p>")

  def test_html_attribute_has(self):
    a, b, c = db.execute("""select 
      html_attribute_has('<p x>', 'p', 'x'), 
      html_attribute_has('<p>abc', 'p', 'x'), 
      html_attr_has('<p> <span z>yo', 'span', 'z' )
    """).fetchone()
    self.assertEqual(a, 1)
    self.assertEqual(b, 0)
    self.assertEqual(c, 1)
  
  def test_html_attribute_get(self):
    a, b, c = db.execute("""select 
      html_attribute_get('<p x>', 'p', 'x'), 
      html_attribute_get('<p>abc', 'p', 'x'), 
      html_attr_get('<p> <span x=z>yo', 'span', 'x' )
    """).fetchone()
    self.assertEqual(a, None)
    self.assertEqual(b, None)
    self.assertEqual(c, "z")
  
  def test_html_extract(self):
    a, b, c = db.execute("""select 
      html_extract('<div> asdfasdf <p a=b>abc</p> asdfasdf </div>', 'p'), 
      html_extract('<p>abc', 'p'), 
      html_extract('<p> <b>asdf</b></p>', 'x')
    """).fetchone()
    self.assertEqual(a, "<p a=\"b\">abc</p>")
    self.assertEqual(b, "<p>abc</p>")
    self.assertEqual(c, None)
  
  def test_html_text(self):
    a, b, c = db.execute("""select 
      html_text('<div> asdfasdf <p a=b>abc</p> asdfasdf </div>', 'p'), 
      html_text('<p>abc', 'p'), 
      html_text('<p> <b>asdf</b></p>', 'x')
    """).fetchone()
    self.assertEqual(a, "abc")
    self.assertEqual(b, "abc")
    self.assertEqual(c, None)

  def test_html_count(self):
    a, b, c = db.execute("""select 
      html_count('<div> ', 'p'), 
      html_count('<div> <p>abc', 'p'), 
      html_count('<div> <p class=x></p> <p>a</p> <p class=x>', 'p.x')
    """).fetchone()
    self.assertEqual(a, 0)
    self.assertEqual(b, 1)
    self.assertEqual(c, 2)
  
  def test_html_each(self):
    rows = db.execute("""select * 
    from html_each('<div>
    <p>a</p>
    <p id=x>b</p>
    <p>c1<span>c2</span></p>
    ', 'p')
    """).fetchall()

    self.assertEqual(len(rows), 3)
    self.assertEqual(len(rows[0]), 3)

    self.assertEqual(rows[0][0], 0)
    self.assertEqual(rows[0][1], "<p>a</p>")
    self.assertEqual(rows[0][2], "a")

    self.assertEqual(rows[1][0], 1)
    self.assertEqual(rows[1][1], "<p id=\"x\">b</p>")
    self.assertEqual(rows[1][2], "b")

    self.assertEqual(rows[2][0], 2)
    self.assertEqual(rows[2][1], "<p>c1<span>c2</span></p>")
    self.assertEqual(rows[2][2], "c1c2")

if __name__ == '__main__':
    unittest.main()