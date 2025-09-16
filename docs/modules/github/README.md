# GitHub

The `github` module exposes functionality for interacting with GitHub issues and pull requests.

Each method exposed by the `github` module is documented below.

## Authentication

By default, the `github` module will attempt to fetch an authentication token
using the GitHub CLI tool (equivalent to `gh auth token -h {hostname}`).

If it is unsuccessful in fetching a token for the GitHub host it is
querying against, it will use anonymous authentication for the requests.

## Methods

### `search`

The `search` method is used to perform a search query using the GitHub search API.

It will return all issues and pull requests that match the provided search query.

This method is subject to limitations outlined in the [GitHub search API documentation](https://docs.github.com/en/rest/search/search?apiVersion=2022-11-28)

As of today, the `search` method does not support pagination and will only return
results from the first page of search results. Additionally, it uses the default
behavior for all additional endpoint parameters present in the
[Search issues and pull requests](https://docs.github.com/en/rest/search/search?apiVersion=2022-11-28#search-issues-and-pull-requests)
query parameter documentation.

Support for specifying these additional query parameters and supporting paginated
results is planned as part of future improvements to this method.


#### Signature

```starlark
github.search(
    host="https://github.com", # Optional. GitHub host to use for API requests. Defaults to github.com.
    query="repo:org/repo is:open label:good-first-issue", # Required. The search query to execute.
    group="good-first-issues", # Optional. A wranglr-specific grouping directive. Useful for conceptual grouping of issues/pull requests.
)
```

#### Return Value

The `search` method will return a Starlark list of all issues and pull requests returned
from the search query execution.

GitHub issues and pull requests are represented like so:
```starlark
items = github.search(...)

item = items[0]

# Get issue/pull request values (immutable)
item.assignees # Get the GitHub handles of assignees. List of strings.
item.author # Get the GitHub handle of the author. String.
item.author_association # Get the association of the author with the project. String.
item.body # Get the body/description of the issue/PR. String.
item.closed_at # Get the datetime the issue/PR was closed. String.
item.comments # Get the number of comments on the issue/PR. Integer.
item.created_at # Get the datetime the issue/PR was created. String.
item.labels # Get the labels present. List of strings.
item.locked # Get whether or not the issue/PR is locked. Boolean.
item.number # Get the number of the issue/PR. Integer.
item.pull_request # Get the pull request data associated with this issue. Dictionary.
item.pull_request["url"] # Get the URL for the pull request. If this is present, the issue is a pull request. String.
item.pull_request["merged_at"] # Get the datetime of when the pull request was merged. String.
item.state # Get the current state of the issue/pull request. String.
item.state_reason # Get the reason for the current state. String.
item.title # Get the title of the issue/pull request. String.
item.updated_at # Get the datetime of the last update. String.

# Get/Set wranglr-specific fields (mutable)
item.status # Represents an arbitrary "status" assigned to this item. Useful in automations for marking things as "Todo", "Needs Review", etc. String.
item.priority # A priority score of the issue. wranglr will sort items in a given view by their priority score. Higher score means higher priority. 64 bit integer.
item.group # Represents a logical "group" this item belongs to. Useful for grouping things into subsets of issues like "Feature X", "SIG Auth", etc.
```

NOTE: datetime values are formatted as RFC-3339 datetime values. Example: `2006-01-02T15:04:05Z07:00`

