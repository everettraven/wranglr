# Jira

The `jira` module exposes functionality for interacting with Jira.

Each method exposed by the `jira` module is documented below.

## Authentication

By default, the `jira` module will attempt to fetch an authentication token
from the environment variable `WRANGLR_JIRA_TOKEN`.

If it is unsuccessful in fetching a token, it will fall back to anonymous
authentication for requests.

It is a known limitation that this restricts usage to fetching data from
a singular Jira host at a time in a given `wranglr` configuration file.

It is planned to add a method for hostname-aware token retrieval in the future
but this limitation should only impact users that use multiple Jira instances
for tracking project work. It is anticipated that this will be a sufficiently
small set of potential users that rolling out the MVP of this implementation
seemed reasonable. If you are impacted by this limitation but are interested
in using `wranglr` let us know you've encountered this limitation through GitHub
discussions/issues so we can better prioritize the improvement.

## Methods

### `search`

The `search` method is used to perform a JQL search query using the Jira search API.

Currently, the `search` method does not support pagination and will only return
the first page of results (~50).

Also worth noting, in testing it has been noticed that searching Jira
has higher latency than when using the `github` module. This is likely
dependent on the instance of Jira you are interacting with, but it is
not unexpected if the use of a handful of expensive Jira search requests
cause it to look like `wranglr` is "hanging".

#### Signature

```starlark
jira.search(
    host="https://issues.host.com", # Required. Jira host to use for API requests.
    query="project = \"Some Project\" AND labels IN (needs-triage)", # Required. The search query to execute.
    group="triaging", # Optional. A wranglr-specific grouping directive. Useful for conceptual grouping of items.
)
```

#### Return Value

The `search` method will return a Starlark list of all tickets returned
from the search query execution.

Jira items are represented like so:
```starlark
items = jira.search(...)

item = items[0]

# Get Jira item values (immutable)
item.assignee # Get the username of the assignee. String or None.
item.creator # Get the username of the author. String or None.
item.reporter # Get the username of the reporter. String or None.
item.type # Get the type of the item (Epic, Story, Bug, etc.). String.
item.project # Get the project this item is associated with. String.
item.resolution # Get the resolution of this item. String or None.
item.ticket_priority # Get the Jira-specific priority of this item. String or None.
item.resolution_date # Get the datetime this item was marked as being resolved. String.
item.created # Get the datetime this item was created. String.
item.due_date # Get the datetime this item is due. String.
item.updated # Get the datetime this item was last updated. String.
item.description # Get the description of this item. String.
item.summary # Get the summary of this item. String.
item.components # Get the components this item is associated with. List of strings.
item.ticket_status # Get the Jira-specific status of this item. String or None.
item.fix_versions # Get the versions this item is "fixed" in. List of strings.
item.affects_versions # Get the versions affected by this item. List of strings.
item.labels # Get the labels present on this item. List of strings.
item.epic # Get the Epic this item belongs to. String or None.
item.sprint # Get the sprint this item is in. String or None.

# Get/Set wranglr-specific fields (mutable)
item.status # Represents an arbitrary "status" assigned to this item. Useful in automations for marking things as "Todo", "Needs Review", etc. String.
item.priority # A priority score of the issue. wranglr will sort items in a given view by their priority score. Higher score means higher priority. 64 bit integer.
item.group # Represents a logical "group" this item belongs to. Useful for grouping things into subsets of issues like "Feature X", "SIG Auth", etc.
```

NOTE: datetime values are formatted as RFC-3339 datetime values. Example: `2006-01-02T15:04:05Z07:00`
