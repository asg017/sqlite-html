import unittest
import subprocess 

class Results:
  def __init__(self, stdout, stderr):
    self.stdout = stdout
    self.stderr = stderr

def run_sqlite3(input):
  if type(input) is list:
    args = ["dist/sqlite3", ":memory:"] + input
  else:
    args = ["dist/sqlite3", ":memory:"] + [input]
  
  proc = subprocess.run(args, stdin=subprocess.PIPE, stdout=subprocess.PIPE)
  out = proc.stdout.decode('utf8') if type(proc.stdout) is bytes else None
  err = proc.stderr.decode('utf8') if type(proc.stderr) is bytes else None
  return Results(out, err)

class TestSqliteLinesCli(unittest.TestCase):
  def test_cli_scalar(self):
    self.assertEqual(run_sqlite3('select 1;').stdout,  '1\n')
    self.assertEqual(
      run_sqlite3(['select name from pragma_function_list where name like "html%" order by 1']).stdout,  
      "html\nhtml_attr_get\nhtml_attr_has\nhtml_attribute_get\nhtml_attribute_has\nhtml_count\nhtml_debug\nhtml_element\nhtml_escape\nhtml_extract\nhtml_table\nhtml_text\nhtml_trim\nhtml_unescape\nhtml_valid\nhtml_version\n"
    )
    self.assertEqual(
      run_sqlite3(['select name from pragma_module_list where name like "html_%" order by 1']).stdout,  
      "html_each\n"
    )
    self.assertEqual(
      run_sqlite3(['select rowid, html, text from html_each("<div> <a>x</a> <a>y</a> <a>z</a>", "a")']).stdout,  
      "0|<a>x</a>|x\n1|<a>y</a>|y\n2|<a>z</a>|z\n"
    )

if __name__ == '__main__':
    unittest.main()