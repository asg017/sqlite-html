# sqlite-dom

## Overview

Scalar functions:

- [dom_all](#dom_all)(_html, selector_)
- [dom_extract](#dom_extract)(_html, selector_)
- [dom_text](#dom_test)(_html, selector_)
- [dom_count](#dom_count)(_html, selector_)
- [dom_trim](#dom_trim)(_text_)
- [dom_table](#dom_-)(_html_)
- [dom_attr_get](#dom_attr_get)(_html, selector, attribute_)

Table-valued functions:

- [dom\_$$](#dom_--)(_html, selector_)

## Interface Overview

selector meaning

## Function Details

### The `dom_$$()` table-vauled function

### The `dom_$()` function

X

Examples:

- `dom_$()` ➡ `''`

### The `dom_$text()` function

X

Examples:

- `dom_$text()` ➡ `''`

### The `dom_trim()` function

X

Examples:

- `dom_trim(' abc ')` ➡ `'abc'`
- `dom_trim('\n\n abc \n\t\t')` ➡ `'abc'`
- `dom_trim('abc')` ➡ `'abc'`

### The `dom_table()` function

X

Examples:

- `dom_table('<tr><td>Alex</td>')` ➡ `'<table> <tr><td>Alex</td></tr>'`

### The `dom_attr_get()` function

X

Examples:

- `dom_attr_get('<a href="https://observablehq.com/@asg017">', 'a', 'href')` ➡ `'https://observablehq.com/@asg017'`

## TODO/Roadmap

### Classes

- [ ] `dom_classes_get(doc, selector)`
- [ ] `dom_class_has(doc, selector, classname)`
- [ ] `dom_class_remove(doc, selector, classname)`

### "Manipulation" Functions

- [ ] `dom_attr_set(doc, selector, name, value)`
- [ ] `dom_attr_remove()`

- [ ] `dom_append(doc, selector, value)`
- [ ] `dom_append_all(doc, selector, value)`
- [ ] `dom_replace(doc, selector, value)`
- [ ] `dom_replace_all(doc, selector, value)`
- [ ] `dom_set(doc, selector, value)`
- [ ] `dom_set_all(doc, selector, value)`
- [ ] `dom_remove(doc, selector)`
- [ ] `dom_remove_all(doc, selector)`

### DOM creation

Only when `subtype` is supported, simiar to json1.

- [ ] `dom_element(name, attrs, ...children)`
- [ ] `dom_attrs(name1, value1, name2, value2, ....)`

