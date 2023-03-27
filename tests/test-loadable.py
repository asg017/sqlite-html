import sqlite3
import unittest

EXT_PATH = "dist/html0"

FUNCTIONS = [
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
    "html_group_element_div",
    "html_group_element_span",
    "html_table",
    "html_text",
    "html_text",
    "html_trim",
    "html_unescape",
    "html_valid",
    "html_version",
  ]
MODULES = [
  "html_each"
]

ALIASES = ["html_attr_get", "html_attr_has"]

def connect(ext):
  db = sqlite3.connect(":memory:")

  db.execute("create table base_functions as select name from pragma_function_list")
  db.execute("create table base_modules as select name from pragma_module_list")

  db.enable_load_extension(True)
  db.load_extension(ext)

  db.execute("create temp table loaded_functions as select name from pragma_function_list where name not in (select name from base_functions) order by name")
  db.execute("create temp table loaded_modules as select name from pragma_module_list where name not in (select name from base_modules) order by name")

  db.row_factory = sqlite3.Row
  return db


db = connect(EXT_PATH)

class TestHtml(unittest.TestCase):
  def test_funcs(self):
    funcs = list(map(lambda a: a[0], db.execute("select name from loaded_functions").fetchall()))
    self.assertEqual(funcs, FUNCTIONS)
  
  def test_modules(self):
    funcs = list(map(lambda a: a[0], db.execute("select name from loaded_modules").fetchall()))
    self.assertEqual(funcs, MODULES)

  def test_html_version(self):
    with open("./VERSION") as f:                                                
      version = "v" + f.read()  
    v, = db.execute("select html_version()").fetchone()
    self.assertEqual(v, version)
  
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
      html_text('<p>abc<p>'), 
      html_text('<p> <b>asdf</b></p>', 'x')
    """).fetchone()
    self.assertEqual(a, "abc")
    self.assertEqual(b, "abc")
    self.assertEqual(c, None)
  
  def test_html_valid(self):
    html_valid = lambda x: db.execute("select html_valid(?)", [x]).fetchone()[0]
    self.assertEqual(html_valid("<div>a"), 1)
    # TODO wtf isn't valid HTML
  
  def test_html_group_element_div(self):
    self.skipTest("")
  def test_html_group_element_span(self):
    self.skipTest("")

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
    rows = db.execute("""select rowid, * 
    from html_each('<div>
    <p>a</p>
    <p id=x>b</p>
    <p>c1<span>c2</span></p>
    ', 'p')
    """).fetchall()

    self.assertEqual(list(map(lambda x: dict(x), rows)), [
      {"rowid":0,"html":"<p>a</p>","text":"a"},
      {"rowid":1,"html":"<p id=\"x\">b</p>","text":"b"},
      {"rowid":2,"html":"<p>c1<span>c2</span></p>","text":"c1c2"}
    ])
    
class TestCoverage(unittest.TestCase):                                      
  def test_coverage(self):                                                      
    test_methods = [method for method in dir(TestHtml) if method.startswith('test_html')]
    funcs_with_tests = set([x.replace("test_", "") for x in test_methods])
    for func in FUNCTIONS:
      if func in ALIASES: continue
      self.assertTrue(func in funcs_with_tests, f"{func} does not have cooresponding test in {funcs_with_tests}")

if __name__ == '__main__':
    unittest.main()