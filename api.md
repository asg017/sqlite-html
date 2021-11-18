# API Reference

As a reminder, `sqlite-html` is still young, so breaking changes should be expected.

## Overview

- Query HTML elements using CSS selectors
  - [html_each](#html_each)(_document, selector_)
  - [html_extract](#html_extract)(_document, selector_)
  - [html_text](#html_text)(_document, selector_)
  - [html_count](#html_count)(_document, selector_)
- Safely generating HTML elements
  - [html](#html)(_document_)
  - [html_element](#html_element)(_tag, attributes, child1, ..._)
- HTML attributes
  - [html_attribute_get](#html_attribute_get)(_document, selector, attribute_)
  - [html_attribute_has](#html_attribute_has)(_document, selector, attribute_)
- Misc. HTML utilities
  - [html_table](#html_table)(_document_)
  - [html_escape](#html_escape)(_text_)
  - [html_unescape](#html_unescape)(_text_)
  - [html_trim](#html_trim)(_text_)
- `sqlite-html` information
  - [html_version](#html_version)()
  - [html_debug](#html_debug)()

## Details

### Query HTML Elements

#### `html_each`

`html_each` is a [table-valued function](https://www.sqlite.org/vtab.html#tabfunc2) that creates a table with the following schema:

```sql
CREATE TABLE html_each(
  i INTEGER,  -- index of the html element, starting at 0
  html TEXT,  -- HTML of the extracted element
  text TEXT,  -- textContent of the HTML element
  length INT  -- length of the HTML element (TODO rm?)
);
```

The `i` column

The `html` column

The `text` column

The `length` column

#### `html_extract`

#### `html_text`

#### `html_count`

### Generate HTML Elements

#### `html`

#### `html_element`

### HTML Attributes

#### `html_attribute_get`

Alias: `html_attr_get`

#### `html_attribute_has`

Alias: `html_attr_has`

### HTML Utilities

#### `html_table`

#### `html_escape`

#### `html_unescape`

#### `html_trim`

Examples:

```sql
sqlite> select html_trim(" asdf ");
asdf
sqlite> select html_trim( html_text("<p>   empty space     </p>", "p") );
empty space
```

### `sqlite-html` Information

#### `html_version`

Returns the version string of the `sqlite-html` library, modeled after [`sqlite_version()`](https://www.sqlite.org/lang_corefunc.html#sqlite_version).

```sql
sqlite> select html_version();
v0.0.0
```

#### `html_debug`

Returns debug information of the `sqlite-html` library, including the version string. Subject to change.

```sql
sqlite> select html_debug();
Version: v0.0.0
Commit: 0cd144a880b47f4a57a5c7f8ceb96eb9dc821508
Runtime: go1.17 darwin/amd64
Date: 2021-11-17T17:06:12Z-0800
```

<h3>The <code><a id="html">html()</a></code> Function </h3>

html(_document_)

Examples:

- `select html("<p class=x>yo");` → `"<p class="x">yo</p>"`
- `select html("<a>");` → `"<a></a>"`

<h3>The <code><a id="html_element">html_element()</a></code> Function </h3>

html*element(*\_tag, attributes, child1, ...\_)

Examples:

```sql
select html_element("p", json_object("class", "greetings"), "hello!");
-- "<p class="greetings">hello!</p>"

select html_element("p", null, "hello! <script></script>")
-- "<p>hello! &lt;script&gt;&lt;/script&gt;</p>"

select html_element("p", null,
  "hello, ",
  html_element("span", null, "Alex"),
  "!"
)
-- '<p>hello, <span>Alex</span>!</p>'

```

<h3>The <code><a id="html_extract">html_extract()</a></code> Function </h3>

html*extract(\_document, selector*)

Examples:

- `html_extract()` → `""`

<h3>The <code><a id="html_text">html_text()</a></code> Function </h3>

html*text(\_document, selector*)

Examples:

- `html_text()` → `""`

<h3>The <code><a id="html_attr_get">html_attr_get()</a></code> Function </h3>

html*attr_get(\_document, selector, attribute*)

Examples:

- `html_attr_get()` → `""`

<h3>The <code><a id="html_attr_has">html_attr_has()</a></code> Function </h3>

html*attr_has(\_document, selector, attribute*)

Examples:

- `html_attr_has()` → `""`

<h3>The <code><a id="html_count">html_count()</a></code> Function </h3>

html*count(\_document, selector*)

Examples:

- `html_count()` → `""`

<h3>The <code><a id="html_table">html_table()</a></code> Function </h3>

html*table(\_document*)

Examples:

- `html_table()` → `""`

<h3>The <code><a id="html_escape">html_escape()</a></code> Function </h3>

html*escape(\_text*)

Examples:

- `html_escape()` → `""`

<h3>The <code><a id="html_unescape">html_unescape()</a></code> Function </h3>

html*unescape(\_text*)

Examples:

- `html_unescape()` → `""`

<h3>The <code><a id="html_trim">html_trim()</a></code> Function </h3>

html*trim(\_text*)

Examples:

- `html_trim()` → `""`

<h3>The <code><a id="html_version">html_version()</a></code> Function </h3>

html_version()

Examples:

- `html_version()` → `""`

<h3>The <code><a id="html_debug">html_debug()</a></code> Function </h3>

html_debug()

Examples:

- `html_debug()` → `""`

## Interface Overview

selector meaning

### The `html_table()` function

X

Examples:

- `html_table('<tr><td>Alex</td>')` ➡ `'<table> <tr><td>Alex</td></tr>'`

### The `html_attr_get()` function

X

Examples:

- `html_attr_get('<a href="https://observablehq.com/@asg017">', 'a', 'href')` ➡ `'https://observablehq.com/@asg017'`
