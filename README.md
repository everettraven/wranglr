# wranglr

`wranglr` is a CLI tool built to reduce the overhead associated with "ticket sprawl"
by creating a unified view of issues, pull requests and tickets across systems
like GitHub and Jira (and more in the future).

It leverages Starlark to provide a scriptable interface so you can choose what
is important, create your own status automation, and prioritize work your way.

![demo](docs/assets/demo.gif)

## Why `wranglr`?

As someone responsible for maintaining multiple projects
it has proven to be a challenge in itself to stay on top
of what is important without losing track of something else.

GitHub notifications and emails aren't sustainable when you get hundreds
a day.

GitHub Projects and Notion have more manual overhead than I'd like.

So, I built `wranglr`.

If you're in a similar boat, maybe you'll find it helpful.

## Installation

```sh
go install github.com/everettraven/wranglr@latest
```

## Quick Start

### Create your `wranglr` configuration file

By default, `wranglr` uses the file `$HOME/.config/wranglr.star`
to determine what sources to fetch from and render.

If that file does not exist it will try finding `wranglr.star` in the current directory.

If you want to use a specific `*.star` file, you can use the `--config` (alias `-c`) flag
to specify the file you'd like `wranglr` to use when running the CLI.

### Writing your `wranglr` configuration file

Here is a succinct example of fetching data from both GitHub and Jira:

```star
# Fetch all open issues and pull requests from the
# kubernetes-sigs/crdify GitHub repository

crdify_open_items = github.search(query="repo:kubernetes-sigs/crdify is:open")

# Fetch all Jira tickets belonging to the project
# "Some Project" with the label "needs-triage"

someproject_untriaged_tickets = jira.search(
    host="https://some.host.com",
    query="project = \"Some Project\" AND label IN (needs-triage)",
)

# Tell wranglr to include everything in the rendered output
wranglr.render(crdify_open_items, someproject_untriaged_tickets)
```

### Running `wranglr`

Once you've written your configuration file, you're ready to run the `wranglr` CLI:
```sh
wranglr
```

This will render the output in a TUI after reading the configuration from `$HOME/.config/wranglr.star`
or the `wranglr.star` file in your current directory.

You can run:
```sh
wranglr --config some/path/config.star
```

to tell `wranglr` that it should use the configuration in the provided file instead.

Don't want the TUI output? You can output JSON as well using:
```sh
wranglr --output json
```


## Contributing

Thanks for your interest in contributing!

The best way to help for now is to:

- Give `wranglr` a try in your workflow
- Share feedback in GitHub discussions / issues
- Report bugs via issues

While you are welcome to submit pull requests, this is a nights/weekends project
which means that responses will be delayed.
