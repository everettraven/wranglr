# synkr

So, you're a developer and your having some trouble keeping track of your work across all the different avenues your company/team uses to track work?

Me too. So I'm building `synkr`. Maybe it will help you as well.

`synkr` is a CLI tool designed to help you fetch "work items" from various sources and filter them via a [Starlark](https://github.com/bazelbuild/starlark) configuration file.

"Work items" is intentionally vague for now as I continue to identify what things developers need to
keep track of in their day-to-day.

Currently, `synkr` only has support for fetching GitHub issues and pull requests.

In the near-term, it is planned to add support for Jira. Longer-term support for other Jira-like tools may be added.

Long-term, the goal of `synkr` is to free engineers of the overhead associated with organizing and keeping track of their tasks allowing them to focus on what they do best - be an engineer.

## Installation

To install `synkr`, download a binary from the releases and add it to your `PATH`.

If you have Go installed, you can also run:

```sh
go install github.com/everettraven/synkr@latest
```

## Usage

### Writing a `synkr` configuration file

#### `synkr` configuration details

`synkr` acts as an engine that processes configurations specified in a Starlark configuration file.

By default, `synkr` will read configuration from `$HOME/.config/synkr.star`. If `synkr` is unable to determine your home directory, it will fallback to `synkr.star` in the current directory.

You can change the file it uses with the `--config` (alias `-c`) flag.

`synkr` has builtin functions that can be used to configure individual sources. Currently there is support for:

- [GitHub](#github)

#### Example `synkr` configuration file

Let's build a quick configuration that allows us to fetch all issues and pull requests
from <https://github.com/kubernetes/kubernetes> where the Kubernetes-SIG API Machinery needs
to provide some input (denoted by the label `sig/api-machinery`):

```starlark
github(repo="kubernetes/kubernetes", labels=["sig/api-machinery"])
```

For more examples, see the `examples/` directory.

### Output

`synkr` currently supports Markdown, JSON and web output formats.

An example of the JSON output (configured with a single source):

```json
{
  "source": "GitHub",
  "project": "kubernetes-sigs/kube-api-linter",
  "items": [
    {
      "id": 3109989899,
      "url": "https://github.com/kubernetes-sigs/kube-api-linter/issues/95",
      "author": "everettraven",
      "labels": [],
      "type": "Issue",
      "assignees": [],
      "title": "Feature: Allow configuration of custom enum markers for `maxlength` linter",
      "body": "In OpenShift, we have some custom markers that set enum values for a field and this results in the `maxlength` linter stating that a field/type alias should have a maximum length when using this custom marker instead of the standard `kubebuilder:validation:Enum` marker.\n\nWhile this particular case is OpenShift-specific, I think it is reasonable to make a generic way to extend this detection logic as there may be other vendors and/or projects that use their own custom markers for CRD generation.",
      "state": "open",
      "priority": 0
    },
    {
      "id": 2590503547,
      "url": "https://github.com/kubernetes-sigs/kube-api-linter/pull/103",
      "author": "everettraven",
      "labels": [
        "cncf-cla: yes",
        "size/M"
      ],
      "type": "PullRequest",
      "assignees": [],
      "title": "markers: fix a bug when parsing expressions with commas present in value",
      "body": "Fixes #99 \r\n\r\nInstead of splitting on solely the `,` character, we now do some more robust normalization for parsing of markers to handle the scenarios where a marker may specify an expression with attributes the have a `,` in their value.",
      "state": "open",
      "priority": 0
    }
  ]
}
```

An example of the Markdown output:

```md
# GitHub - kubernetes-sigs/kube-api-linter
## [Issue][open]: Feature: Allow configuration of custom enum markers for `maxlength` linter
**URL**: https://github.com/kubernetes-sigs/kube-api-linter/issues/95
**Author**: *everettraven*
**Assignees**:



In OpenShift, we have some custom markers that set enum values for a field and this results in the `maxlength` linter stating that a field/type alias should have a maximum length when using this custom marker instead of the standard `kubebuilder:validation:Enum` marker.

While this particular case is OpenShift-specific, I think it is reasonable to make a generic way to extend this detection logic as there may be other vendors and/or projects that use their own custom markers for CRD generation.

## [PullRequest][open]: markers: fix a bug when parsing expressions with commas present in value
**URL**: https://github.com/kubernetes-sigs/kube-api-linter/pull/103
**Author**: *everettraven*
**Assignees**:

`cncf-cla: yes` `size/M`

Fixes #99

Instead of splitting on solely the `,` character, we now do some more robust normalization for parsing of markers to handle the scenarios where a marker may specify an expression with attributes the have a `,` in their value.
```

The web output starts an HTTP server on an available localhost port to serve a single page application that shows the data in a paginated table that can be searched and sorted.

### Sources

#### GitHub

In order to use the GitHub source, you use the `github` builtin function like so:

```starlark
github(
    host?="github.com",
    repo="{organization}/{repository}",
    keywords?=[""],
    assignee?="",
    author?="",
    closed?="YYYY-MM-DD",
    commenter?="",
    comments?="",
    created?="YYYY-MM-DD",
    draft?="true | false",
    extension?="",
    filename?="",
    in?=[""],
    involves?="",
    is?=[""],
    labels?=[""],
    language?="",
    mentions?="",
    merged?="YYYY-MM-DD",
    milestone?="",
    no?=[""],
    path?="",	
    review?="",
	review_requested?="",
	reviewed_by?="",
	state?="",
	team?="",
	team_review_requested?="",
	updated?="YYYY-MM-DD",
	sort?="",
	order?="",
	limit?=100,
    filters?=[...],
    priorities?=[...],
    status?={function},
)
```

##### Request Filter parameters

The following describes how each of the parameters in the `github` builtin function influences the request
made against the GitHub API.

- `host` is the hostname for the GitHub instance to query the API of. Optional. Defaults to `github.com`.
- `repo` is the organization and repository pair to pull issues and pull requests from. Required.
- `keywords` is the set of keywords to filter search results on. Optional.
- `assignee` requests only issues and pull request assigned to this particular user. Optional.
- `author` requests only issues and pull requests authored by this particular user. Optional.
- `closed` requests only issues and pull requests where they have been closed in the provided date range (i.e `>2024-12-1` ,`<2025-1-1`, `2024-12-1..2025-1-1`). Optional.
- `commenter` requests only issues and pull requests the specified user has commented on. Optional.
- `comments` requests only issues and pull requests where the number of comments matches the specified constraint (i.e `>10`, `<20`, `10..20`). Optional.
- `created` requests only issues and pull requests where the creation date is in the provided date range. Optional.
- `draft` requests only pull requests that match the draft constraint (`true` or `false`). Optional.
- `extension` requests only pull requests that contain changes to files with the specified file extension. Optional.
- `filename` requests only pull requests that contain changes to the specified filename. Optional.
- `in` requests only pull requests where the keywords exist in the specified fields (i.e `title`, `body`). Optional.
- `involves` requests only issues and pull requests that involve the specified user. Optional.
- `is` requests only issues and pull requests that match the specified constraints (`issue`, `pr`, `draft`, `open`, `closed`, `merged`). Optional.
- `labels` requests only issues and pull requests that have the specified set of labels. Optional.
- `language` requests only issues and pull requests that contain the specified programming language. Optional.
- `mentions` requests only issues and pull requests where the specified user has been mentioned. Optional.
- `merged` requests only pull requests that have merged based on the date criteria set. Optional.
- `milestone` requests only issues and pull requests that are assigned to the specified milestone. Optional.
- `no` requests only issues and pull requests that do not have the specified attributes (i.e `milestone`, `assignee`, `labels`). Optional.
- `path` requests only pull requests that contain changes to the specified path. Optional.
- `review` requests only pull requests that have the specified review state (i.e `none`, `approved`, `changes_requested`). Optional.
- `review_requested` requests only pull requests that have requested a review from the specified user. Optional.
- `reviewed_by` requests only pull requests that have been reviewed by the specified user. Optional.
- `state` requests only issues and pull requests based on the specified state (i.e `open`, `closed`). Optional.
- `team` requests only issues and pull requests where the specified team is mentioned. Optional.
- `team_review_requested` requests only pull requests where the specified team has been requested for review. Optional.
- `updated` requests only issues and pull requests that have been updated based on the date criteria specified. Optional.
- `sort` specifies how the returned results from the GitHub search API should be sorted (i.e `created`, `updated`, `comments`, `reactions`). Optional.
- `order` specifies the ordering in which results from the GitHub search API should be returned (i.e `asc`, `desc`). Optional.
- `limit` specifies the maximum number of results that should be returned from the GitHub search API. Optional.

##### Post-Request Parameters

The following describes how each parameter influences the processing of the results returned by the GitHub search API: 

`filters` is an optional list of functions that should be called by `synkr` when determining whether or not an issue or pull request should be included in the returned set.
The functions are expected to accept a single parameter and return a "truthy" value (i.e `True` / `False` state should be able to be determined from the returned value).
A return value reflective of the `True` state means that an item should be included in the output.
The parameter passed to the functions is a dictionary with the following keys and value types:

- `author` (string). The GitHub handle of the author of the issue/pull request. Example: `everettraven`.
- `type` (string). The type of item this is, one of `Issue`, `PullRequest`.
- `title` (string). The title of the issue/pull request.
- `body` (string). The body of the issue/pull request.
- `state` (string). The current state of the issue/pull request. Example: `open`/`closed`.
- `labels` ([]string). The current set of labels on the issue/pull request.
- `assignees` ([]string). The current set of assignees on the issue/pull request.
- `created` (string). The creation timestamp of the issue/pull request. Example: `2025-04-10 18:34:13 +0000 UTC`
- `updated` (string). The timestamp of the last time the issue/pull request was updated.
- `comments` (int). The number of comments on the issue/pull request.
- `milestone` (string). The title of the milestone the issue/pull request is included in.
- `mentions` ([]string). The GitHub handles of users that are mentioned in the issue/pull request comments. Only populated if `include_mentions` is set to `True` in the `github` builtin function.
- `requestedReviewers` ([]string). The GitHub handles of users whose reviews were explicitly requested on a pull request. Only populated on items where `type` is `PullRequest`.
- `draft` (boolean). Whether or not the pull request is a draft. Only populated on items where `type` is `PullRequest`

`priorities` is an optional list of functions that should be called by `synkr` when determining the priority score to assign to an issue or pull request.
The functions are expected to accept a single parameter (the same parameter as `filters` functions) and return an integer value to add to the item's priority score.

`status` is an optional function that should be called by `synkr` when determining the "status" to assign to an issue or pull request.
"status" is a distinctly different value than `state`, as it represents an arbitrary status defined by you the user instead of GitHub's perceived state of the item.
The function is expected to accept a single parameter (the same parameter as above) and return a string value to set the item's status to.

##### Authentication

By default, the GitHub source will use the unauthenticated GitHub API to fetch issues and pull requests from the configured repositories. This means you will only be able to access public repositories.

In order to access private repositories, you can [create a fine-grained personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token) and set the `SYNKR_GITHUB_TOKEN` environment variable with this token.

### Help

```sh
  synkr is an engine for syncing work items based on a Starlark configuration

  USAGE


    synkr [command] [--flags]  


  COMMANDS

    completion [command]  Generate the autocompletion script for the specified shell
    help [command]        Help about any command

  FLAGS

     -c --config          Configures the Starlark file to be processed for configuration. Defaults to $HOME/.config/synkr.star if possible to get your home directory. Otherwise it uses synkr.star in the current directory. (synkr.star)
     -h --help            Help for synkr
     -o --output          Configures the output format. Allowed values are [markdown, json, web] (markdown)
     -v --version         Version for synkr
```

## Contributing

Thanks for your interest in contributing!

The most impactful contribution today would be to take `synkr` for a spin
and share your thoughts. Please feel free to use GitHub discussions for sharing your
thoughts.

Something broken? Feel free to submit an issue and I'll take a look as soon as I can.

Want to contribute some code? Go for it! I'm open to accepting code contributions pending they align
with the project direction or solve an existing issue that has been discussed and determined warrants a fix.
