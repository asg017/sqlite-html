# API Reference

As a reminder, `sqlite-html` is still young, so breaking changes should be expected while `sqlite-html` is in a pre-v1 stage.

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

A [table function](https://www.sqlite.org/vtab.html#tabfunc2) with the following schema:

```sql
CREATE TABLE html_each(
  i INTEGER,  -- index of the html element, starting at 0
  html TEXT,  -- HTML of the extracted element
  text TEXT  -- textContent of the HTML element
);
```

The `i` column contains the index (starts at 0) of the matching element.

The `html` column contains the matching element's HTML representation.

The `text` column contains the matching element's textContent representation, similar to the JavaScript DOM API's `.textContent` or the `html_text` function in this library.

```sql
sqlite> select * from html_each('<ul>
<li>Alpha</li>
<li>Bravo</li>
<li>Charlie</li>
<li>Delta</li>', 'li')

```

#### `html_extract`

`html_extract(document, selector)`

Extracts the first matching element from `document` using the given CSS `selector`, and returns the full HTML representation of that element.

```sql
sqlite> select html_extract("<p> Hello, <b class=x>world!</b> </p>", "b");
<b class="x">world!</b>

```

#### `html_text`

`html_text(document, selector)`

Extracts the first matching element from `document` using the given CSS `selector`, and returns the text representation of that element, Similar to the [`Node.textContent`](https://developer.mozilla.org/en-US/docs/Web/API/Node/textContent) property in the JavaScript DOM API.

Examples:

```sql
sqlite> select
```

#### `html_count`

`html_count(document, selector)`

For the given `document`, count the number of matching elements from `selector` and return that number.

```sql
sqlite> select
```

### Generate HTML Elements

#### `html`

```sql
sqlite> select
```

#### `html_element`

```sql
sqlite> select
```

### HTML Attributes

#### `html_attribute_get`

Alias: `html_attr_get`

```sql
sqlite> select
```

#### `html_attribute_has`

Alias: `html_attr_has`

```sql
sqlite> select
```

### HTML Utilities

#### `html_table`

```sql
sqlite> select
```

#### `html_escape`

```sql
sqlite> select
```

#### `html_unescape`

```sql
sqlite> select
```

#### `html_trim`

```sql
sqlite> select html_trim(" asdf ");
asdf
sqlite> select html_trim( html_text("<p>   empty space     </p>", "p") );
empty space
```

### `sqlite-html` Information

#### `html_version`

Returns the version string of the `sqlite-html` library, modeled after [`sqlite_version()`](https://www.sqlite.org/lang_corefunc.html#sqlite_version).

Examples:

```sql
sqlite> select html_version();
v0.0.0
```

#### `html_debug`

Returns debug information of the `sqlite-html` library, including the version string. Subject to change.

Examples:

```sql
sqlite> select html_debug();
Version: v0.0.0
Commit: 0cd144a880b47f4a57a5c7f8ceb96eb9dc821508
Runtime: go1.17 darwin/amd64
Date: 2021-11-17T17:06:12Z-0800
```
