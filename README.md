# sqlite-html

A SQLite extension for querying, manipulating, and creating HTML elements.

- Extract HTML or text from HTML with CSS selectors, like `.querySelector()`, `.innerHTML`, and `.innerText`
- Generate a table of matching elements from a CSS selector, like `.querySelectorAll()`
- Safely create HTML elements in a query, like `.createElement()` and `.appendChild()`

`sqlite-html`'s API is modeled after the official [JSON1](https://www.sqlite.org/json1.html#jmini) SQLite extension.

This extension is written in Go, thanks to [riyaz-ali/sqlite](https://github.com/riyaz-ali/sqlite). While this library aims to be fast and efficient, it is overall slower than what a pure C SQLite extension could be, but in practice you may not notice much of a difference.

## Usage

```sql
.load ./html0
select html_extract('<p> Anakin <b>Skywalker</b> </p>', 'b');
-- "<b>Skywalker</b>"
```

`sqlite-html` is similar to other HTML scraping tools like [BeautifulSoup](https://beautiful-soup-4.readthedocs.io/en/latest/) (Python) or [cheerio](https://cheerio.js.org/) (Node.js) or [nokogiri](https://nokogiri.org/) (Ruby). You can use CSS selectors to extract individual elements or groups of elements to query data from HTML sources.

For example, here we find all `href` links in an `index.html` file.

```sql
select
  text as name,
  html_attribute_get(anchors, 'a', 'href') as href
from html_each(readfile('index.html'), 'a') as anchors
```

We can also safely generate HTML with `html_element`, modeled after React's [`React.createElement`](https://reactjs.org/docs/react-api.html#createelement).

```sql
select html_element('p', null,
  'Luke, I am your',
  html_element('b', null, 'father'),
  '!',

  html_element('img', json_object(
    'src', 'https://images.dog.ceo/breeds/groenendael/n02105056_4600.jpg',
    'width', 200
  ))
);

-- "<p>Luke, I am your<b>father</b>!<img src="https://images.dog.ceo/breeds/groenendael/n02105056_4600.jpg" width="200.000000"/></p>"
```

## Documentation

See [`docs.md`](./docs.md) for a full API reference.

## Installing

| Language       | Install                                                      |                                                                                                                                                                                             |
| -------------- | ------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Python         | `pip install sqlite-html`                                    | [![PyPI](https://img.shields.io/pypi/v/sqlite-html.svg?color=blue&logo=python&logoColor=white)](https://pypi.org/project/sqlite-html/)                                                      |
| Datasette      | `datasette install datasette-sqlite-html`                    | [![Datasette](https://img.shields.io/pypi/v/datasette-sqlite-html.svg?color=B6B6D9&label=Datasette+plugin&logoColor=white&logo=python)](https://datasette.io/plugins/datasette-sqlite-html) |
| Node.js        | `npm install sqlite-html`                                    | [![npm](https://img.shields.io/npm/v/sqlite-html.svg?color=green&logo=nodedotjs&logoColor=white)](https://www.npmjs.com/package/sqlite-html)                                                |
| Deno           | [`deno.land/x/sqlite_html`](https://deno.land/x/sqlite_html) | [![deno.land/x release](https://img.shields.io/github/v/release/asg017/sqlite-html?color=fef8d2&include_prereleases&label=deno.land%2Fx&logo=deno)](https://deno.land/x/sqlite_html)        |
| Ruby           | `gem install sqlite-html`                                    | ![Gem](https://img.shields.io/gem/v/sqlite-html?color=red&logo=rubygems&logoColor=white)                                                                                                    |
| Github Release |                                                              | ![GitHub tag (latest SemVer pre-release)](https://img.shields.io/github/v/tag/asg017/sqlite-html?color=lightgrey&include_prereleases&label=Github+release&logo=github)                      |

<!--
| Elixir         | [`hex.pm/packages/sqlite_html`](https://hex.pm/packages/sqlite_html) | [![Hex.pm](https://img.shields.io/hexpm/v/sqlite_html?color=purple&logo=elixir)](https://hex.pm/packages/sqlite_html)                                                                       |
| Go             | `go get -u github.com/asg017/sqlite-html/bindings/go`               | [![Go Reference](https://pkg.go.dev/badge/github.com/asg017/sqlite-html/bindings/go.svg)](https://pkg.go.dev/github.com/asg017/sqlite-html/bindings/go)                                     |
| Rust           | `cargo add sqlite-html`                                             | [![Crates.io](https://img.shields.io/crates/v/sqlite-html?logo=rust)](https://crates.io/crates/sqlite-html)                                                                                 |
-->

The [Releases page](https://github.com/asg017/sqlite-html/releases) contains pre-built binaries for Linux amd64, MacOS amd64 (no arm), and Windows.

### As a loadable extension

If you want to use `sqlite-html` as a [Runtime-loadable extension](https://www.sqlite.org/loadext.html), Download the `html0.dylib` (for MacOS), `html0.so` (Linux), or `html0.dll` (Windows) file from a release and load it into your SQLite environment.

> **Note:**
> The `0` in the filename (`html0.dylib`/ `html0.so`/`html0.dll`) denotes the major version of `sqlite-html`. Currently `sqlite-html` is pre v1, so expect breaking changes in future versions.

For example, if you are using the [SQLite CLI](https://www.sqlite.org/cli.html), you can load the library like so:

```sql
.load ./html0
select html_version();
-- v0.0.1
```

Or in Python, using the builtin [sqlite3 module](https://docs.python.org/3/library/sqlite3.html):

```python
import sqlite3

con = sqlite3.connect(":memory:")

con.enable_load_extension(True)
con.load_extension("./html0")

print(con.execute("select html_version()").fetchone())
# ('v0.0.1',)
```

Or in Node.js using [better-sqlite3](https://github.com/WiseLibs/better-sqlite3):

```javascript
const Database = require("better-sqlite3");
const db = new Database(":memory:");

db.loadExtension("./html0");

console.log(db.prepare("select html_version()").get());
// { 'html_version()': 'v0.0.1' }
```

Or with [Datasette](https://datasette.io/):

```
datasette data.db --load-extension ./html0
```

## See also

- [sqlite-http](https://github.com/asg017/sqlite-http), for making HTTP requests in SQLite (pairs great with this tool)
- [htmlq](https://github.com/mgdm/htmlq), for a similar but CLI-based HTML query tool using CSS selectors
- [riyaz-ali/sqlite](https://github.com/riyaz-ali/sqlite), the brilliant Go library that this library depends on
- [nalgeon/sqlean](https://github.com/nalgeon/sqlean), several pre-compiled handy SQLite functions, in C
