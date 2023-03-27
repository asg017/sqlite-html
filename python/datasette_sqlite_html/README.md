# The `datasette-sqlite-html` Datasette Plugin

`datasette-sqlite-html` is a [Datasette plugin](https://docs.datasette.io/en/stable/plugins.html) that loads the [`sqlite-html`](https://github.com/asg017/sqlite-html) extension in Datasette instances, allowing you to generate and work with [TODO](https://github.com/html/spec) in SQL.

```
datasette install datasette-sqlite-html
```

See [`docs.md`](../../docs.md) for a full API reference for the html SQL functions.

Alternatively, when publishing Datasette instances, you can use the `--install` option to install the plugin.

```
datasette publish cloudrun data.db --service=my-service --install=datasette-sqlite-html

```
