# wranglr

The `wranglr` module exposes functionality specific to `wranglr`

Each method exposed by the `wranglr` module is documented below.

## Methods

### `render`

The `render` method is used to tell `wranglr` which items to include in the rendered output.

The output format that is used for rendering is determined by the `--output`
(alias `-o`, defaults to `interactive`) flag when running the `wranglr` binary.

#### Signature

```starlark
wranglr.render([item, item2], [item, item2], ...) # Lists of items to be rendered.
```

#### Return Value

The `render` method has no return value.

It is specifically used to tell `wranglr` to display the provided items in the desired output format.

Specifying multiple `wranglr.render(...)` calls in a single configuration will result in multiple
render steps executing in succession. For example, a configuration like:

```starlark
items = github.search(...)
wranglr.render(items)

more_items = github.search(...)
wranglr.render(more_items)
```

Will perform a blocking render on the first call.
Once that rendering process is complete, the configuration will continue
to be executed until the next render call.
