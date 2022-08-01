# API Reference

As a reminder, `sqlite-html` is still young, so breaking changes should be expected while `sqlite-html` is in a pre-v1 stage.

## Overview

- `sqlite-html` information
  - [html_version](#html_version)()
  - [html_debug](#html_debug)()
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
  - [html_escape](#html_escape)(_text_)
  - [html_unescape](#html_unescape)(_text_)
  - [html_trim](#html_trim)(_text_)
  - [html_table](#html_table)(_document_)

### Query HTML Elements

#### `html_each()`

A [table function](https://www.sqlite.org/vtab.html#tabfunc2) with the following schema:

```sql
CREATE TABLE html_each(
  html TEXT,  -- HTML of the extracted element
  text TEXT,  -- textContent of the HTML element

  document TEXT hidden, -- input HTML document
  selector TEXT hidden -- input CSS selector
);
```

The `html` column contains the matching element's HTML representation.

The `text` column contains the matching element's textContent representation, similar to the JavaScript DOM API's `.textContent` or the `html_text` function in this library.

```sql
sqlite> select * from html_each('<ul>
<li>Alpha</li>
<li>Bravo</li>
<li>Charlie</li>
<li>Delta</li>', 'li')

```

#### `html_extract(document, selector)`

Extracts the first matching element from `document` using the given CSS `selector`, and returns the full HTML representation of that element.

```sql
select html_extract('<p> Hello, <b class=x>world!</b> </p>', 'b');
-- '<b class="x">world!</b>'

```

#### `html_text(document, selector)`

Extracts the first matching element from `document` using the given CSS `selector`, and returns the text representation of that element, Similar to the [`Node.textContent`](https://developer.mozilla.org/en-US/docs/Web/API/Node/textContent) property in the JavaScript DOM API.

Examples:

```sql
select html_text('<p> hello <a href="https://google.com">dog</a></a>', 'a');
-- "dog"
```

#### `html_count(document, selector)`

For the given `document`, count the number of matching elements from `selector` and return that number.

```sql
select html_count('<div> <p>a</p> <p>b</p> <p>c</p> </div>', 'p');
-- 3
```

### Generate HTML Elements

#### `html(contents)`

Verifies and "cleans" (quotes attributes) the given document as HTML.

Also sets the return subtype to the HTML magic number, for use in other funcs like html_element to designate something as "HTML".

```sql

select html('<a> foo'); -- "<a> foo</a>"

-- returns specific subtype to denote that it's HTML
select subtype('alex'); -- 0
select subtype(html('alex')); -- 221
```

#### `html_element(type, [attributes], [...children])`

```sql

-- Tag only

select html_element('a'); -- "<a></a>"
select html_element('br'); -- "<br/>"

-- attributes passed in as JSON

select html_element('img', json_object('src', './a.png', 'width', 200)); -- '<img src="./a.png" width="200.000000"/> '

-- Children can be text or HTML

select html_element('p', null, "text node"); -- '<p>text node</p>'

select html_element('p', null, "<b>Still a text node</b>"); -- '<p>&lt;b&gt;Still a text node&lt;/b&gt;</p>'

select html_element('p', null, html_element('b', null, 'Bolded!')); -- '<p><b>Bolded!</b></p>'

select html_element('p', null, html('<b>Also bolded</b>')); -- '<p><b>Also bolded</b></p>'

select html_element('p', null,
  "multiple ",
  html("<b>children"),
  " works ",
  html_element("span", null, "just fine")
); -- '<p>multiple <b>children</b>works <span>just fine</span></p>'

```

### HTML Attributes

#### `html_attribute_get(document, selector, attribute)`

Get the value of the "name" attribute from the element found in document, using selector

Alias: `html_attr_get`

```sql
select html_attr_get('<p> <a href="./about"> About<a/> </p>', 'a', 'href'); -- './about'

select html_attr_get('<p> <a href="./about"> About<a/> </p>', 'a', 'rel'); -- NULL
```

#### `html_attribute_has(document, selector, attribute)`

Returns 1 or 0, if the "name" attribute from the element found in document, using selector, exists.

Alias: `html_attr_has`

```sql
select html_attr_has('<p> <a href="./about"> About<a/> </p>', 'a', 'href'); -- 1

select html_attr_has('<p> <a href="./about"> About<a/> </p>', 'a', 'rel'); -- 0
```

### HTML Utilities

#### `html_escape(content)`

Returns an HTML escaped version of the given content.

```sql
select html_escape('<a>');
-- "&lt;a&gt;"
```

#### `html_unescape(content)`

```sql
select html_unescape('&lt;a');
-- "<a"
```

#### `html_trim(contents)`

Trims whitespace around `contents`. Useful since many results of `html_text` will have newlines/spaces that aren't useful.

```sql
select html_trim(" asdf ");
-- "asdf"

select html_trim( html_text("<p>   empty space     </p>", "p") );
-- "empty space"
```

#### `html_table(contents)`

Prepend the string `"<table>"` before `contents`.

While seemingly unnecessary, it's useful when dealing with HTML table `<tr>` rows. The string `<tr>example <a>foo</a> </tr>` isn't valid HTML, so trying to extract `"foo"` won't work by itself. So wrapping `html_table()` before extracting from a `tr` element will allow for `sqlite-html` to parse correctly.

```sql
select html_table('xyz');
-- "<table>xyz"

-- Try removing the "html_table()" calls and see blank results
select
  html_text(html_table(rows.html), 'td:nth-child(1)') as name,
  html_text(html_table(rows.html), 'td:nth-child(2)') as age
from html_each('<table>
<tr> <td>Alex</td> <td>1</td> </tr>
<tr> <td>Brian</td> <td>2</td> </tr>
<tr> <td>Craig</td> <td>3</td> </tr>
</table>', 'tr') as rows;
/*
┌───────┬─────┐
│ name  │ age │
├───────┼─────┤
│ Alex  │ 1   │
│ Brian │ 2   │
│ Craig │ 3   │
└───────┴─────┘
*/
```

### `sqlite-html` Information

#### `html_version()`

Returns the version string of the `sqlite-html` library, modeled after [`sqlite_version()`](https://www.sqlite.org/lang_corefunc.html#sqlite_version).

Examples:

```sql
sqlite> select html_version();
v0.0.0
```

#### `html_debug()`

Returns debug information of the `sqlite-html` library, including the version string. Subject to change.

Examples:

```sql
sqlite> select html_debug();
Version: v0.0.0
Commit: 0cd144a880b47f4a57a5c7f8ceb96eb9dc821508
Runtime: go1.17 darwin/amd64
Date: 2021-11-17T17:06:12Z-0800
```
