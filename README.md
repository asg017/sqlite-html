# sqlite-html

A SQLite extension for querying, manipulating, and creating HTML elements.

- Extract HTML or text from HTML with CSS selectors, like `.querySelector()`, `.innerHTML`, and `.innerText`
- Generate a table of matching elements from a CSS selector, like `.querySelectorAll()`
- Safely create HTML elements in a query, like `.createElement()` and `.appendChild()`

## ⚠️⚠️WARNING⚠️⚠️

This extension is very new and in beta! The API is not stable and should be used with caution. Eventually I'll come up with a roadmap for a stable `v1`, but for now expect breakages in the future.

## Installing

TODO

## Overview

`sqlite-html`'s API is modeled after the official [JSON1](https://www.sqlite.org/json1.html#jmini) SQLite extension. This extension is also written in Go, thanks to [riyaz-ali/sqlite](https://github.com/riyaz-ali/sqlite). While this library aims to be fast and efficient, it is overall slower than what a pure C SQLite extension could be (mostly because cgo is slow), but in practice you probably won't notice much of a difference.

## API Reference

Scalar functions:

- [html](#html)(_document_)
- [html_element](#html_element)(_tag, attributes, child1, ..._)
- [html_extract](#html_extract)(_document, selector_)
- [html_text](#html_text)(_document, selector_)
- [html_attr_get](#html_attr_get)(_document, selector, attribute_)
- [html_attr_has](#html_attr_has)(_document, selector, attribute_)
- [html_count](#html_count)(_document, selector_)
- [html_table](#html_table)(_document_)
- [html_escape](#html_escape)(_text_)
- [html_unescape](#html_unescape)(_text_)
- [html_trim](#html_trim)(_text_)
- [html_version](#html_version)()
- [html_debug](#html_debug)()

Table functions:

- [html_each](#html_each)(_document, selector_)

## Functions

<h3>The <code><a id="html">html()</a></code> Function </h3>

html(_document_)

Examples:

- `select html("<p class=x>yo");` → `"<p class="x">yo</p>"`
- `select html("<a>");` → `"<a></a>"`

<h3>The <code><a id="html_element">html_element()</a></code> Function </h3>

html_element(_\_tag, attributes, child1, ..._)

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

## Function Details

### The `html_$$()` table-vauled function

### The `html_$()` function

X

Examples:

- `html_$()` ➡ `''`

### The `html_$text()` function

X

Examples:

- `html_$text()` ➡ `''`

### The `html_trim()` function

X

Examples:

- `html_trim(' abc ')` ➡ `'abc'`
- `html_trim('\n\n abc \n\t\t')` ➡ `'abc'`
- `html_trim('abc')` ➡ `'abc'`

### The `html_table()` function

X

Examples:

- `html_table('<tr><td>Alex</td>')` ➡ `'<table> <tr><td>Alex</td></tr>'`

### The `html_attr_get()` function

X

Examples:

- `html_attr_get('<a href="https://observablehq.com/@asg017">', 'a', 'href')` ➡ `'https://observablehq.com/@asg017'`
