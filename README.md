# synkr

So, you're a developer and your having some trouble keeping track of your work across all the different avenues your company/team uses to track work?

Me too. So I'm building `synkr`. Maybe it will help you as well.

`synkr` is a CLI tool designed to help you fetch "work items" from various sources and filter them via a [Starlark](https://github.com/bazelbuild/starlark) configuration file.

"Work items" is intentionally vague for now as we continue to identify what things developers need to
keep track of in their day-to-day.

Currently, `synkr` only has support for fetching public GitHub issues and pull requests.

In the future, `synkr` might support fetching items from things like Jira, Google Docs, Slack, and beyond.

Long-term, the goal of `synkr` is to free engineers of the overhead associated with organizing and keeping track of their tasks allowing them to focus on what they do best - engineer.

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

By default, `synkr` will read a `synkr.star` file in the current directory. You can change the file it uses with the `--config` (alias `-c`) flag.

`synkr` has a builtin for configuring individual GitHub sources like so:
```starlark
github(org="kubernetes", repo="kubernetes", filters=[...])
```
`filters` is an optional list of functions that follow the pattern:
```starlark
def filter(item):
  ...
```
where `item` is a dictionary. An example of an item passed to the function:
```json
{
  "id": 3109989899,
  "url": "https://github.com/kubernetes-sigs/kube-api-linter/issues/95",
  "author": "everettraven",
  "labels": [],
  "type": "Issue",
  "assignees": [],
  "title": "Feature: Allow configuration of custom enum markers for `maxlength` linter",
  "body": "In OpenShift, we have some custom markers that set enum values for a field and this results in the `maxlength` linter stating that a field/type alias should have a maximum length when using this custom marker instead of the standard `kubebuilder:validation:Enum` marker.\n\nWhile this particular case is OpenShift-specific, I think it is reasonable to make a generic way to extend this detection logic as there may be other vendors and/or projects that use their own custom markers for CRD generation.",
  "state": "open"
}
```

The above example item does not currently have any labels or assignees, but the respective fields would be populated in the item with the names of labels on the issue and the GitHub handles of assignees respectively for issues/pull requests that do have these fields populated.

#### Example `synkr` configuration file

Let's build a quick configuration that allows us to fetch all issues and pull requests
from https://github.com/kubernetes/kubernetes where the Kubernetes-SIG API Machinery needs
to provide some input (denoted by the label `sig/api-machinery`):

```starlark
def has_sig_api_machinery(item):
  labels = item.get("labels")
  if "sig/api-machinery" in labels:
    return True
  return False

github(org="kubernetes", repo="kubernetes", filters=[has_sig_api_machinery])
```
For more examples, see the `examples/` directory.

### Output

`synkr` currently supports Markdown and JSON output formats.

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
      "state": "open"
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
      "state": "open"
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

### Help
```sh

  synkr is an engine for syncing work items based on a Starlark configuration

  USAGE


    synkr [command] [-c configFile] [--flags]  


  COMMANDS

    completion [command]      Generate the autocompletion script for the specified shell
    help [command] [--flags]  Help about any command

  FLAGS

     -c --config              Configures the Starlark file to be processed for configuration (synkr.star)
     -h --help                Help for synkr
     -v --version             Version for synkr

```

## Contributing

Thanks for your interest in contributing!

The most impactful contribution today would be to take `synkr` for a spin
and share your thoughts. Please feel free to use GitHub discussions for sharing your
thoughts.

Something broken? Feel free to submit an issue and I'll take a look as soon as I can.

Want to contribute some code? Go for it! I'm open to accepting code contributions pending they align
with the project direction or solve an existing issue that has been discussed and determined warrants a fix.
