# sqlite-html

A SQLite extension for querying, manipulating, and creating HTML elements.

- Extract HTML or text from HTML with CSS selectors, like `.querySelector()`, `.innerHTML`, and `.innerText`
- Generate a table of matching elements from a CSS selector, like `.querySelectorAll()`
- Safely create HTML elements in a query, like `.createElement()` and `.appendChild()`

## 🚧🚧 Work In Progress! 🚧🚧

This library is experimental and subject to change. I plan to make a stable beta release and subsequent v0+v1 in the near future, so use with caution.

When v0 is ready (with a mostly stable API), I will make a release (so watch this repo for that) and will make a blog post, feel free to [follow me on twitter](https://twitter.com/agarcia_me) to get notified of that.

## Installing

TODO

## Documentation

See [`api.md`](./api.md) for a full API reference.

## Overview

`sqlite-html`'s API is modeled after the official [JSON1](https://www.sqlite.org/json1.html#jmini) SQLite extension.

This extension is also written in Go, thanks to [riyaz-ali/sqlite](https://github.com/riyaz-ali/sqlite). While this library aims to be fast and efficient, it is overall slower than what a pure C SQLite extension could be (mostly because cgo is slow), but in practice you probably won't notice much of a difference.

## See also

- [sqlite-http](https://github.com/asg017/sqlite-http), for making HTTP requests in SQLite (pairs great with this tool)
- [htmlq](https://github.com/mgdm/htmlq), for a similar but CLI-based HTML query tool using CSS selectors
- [riyaz-ali/sqlite](https://github.com/riyaz-ali/sqlite), the brilliant Go library that this library depends on
- [nalgeon/sqlean](https://github.com/nalgeon/sqlean), several pre-compiled handy SQLite functions, in C
