# Example Configurations

## Open Pull Requests I Authored

```starlark
prs = github.search(query="is:pr is:open author:@me")

wranglr.render(prs)
```

## Kubernetes Issues grouped by SIG

```starlark
sig_auth_issues = github.search(query="repo:kubernetes/kubernetes is:open is:issue label:sig/auth", group="SIG Auth")

sig_apimachinery_issues = github.search(query="repo:kubernetes/kubernetes is:open is:issue label:sig/api-machinery", group="SIG API Machinery")

wranglr.render(sig_auth_issues, sig_apimachinery_issues)
```

## Setting Status of Items

```starlark
kubernetes_issues = github.search(query="repo:kubernetes/kubernetes is:open is:issue")

for issue in kubernetes_issues:
  if "triage/accepted" not in issue.labels:
    issue.status = "Needs Triage"

wranglr.render(kubernetes_issues)
```

## Prioritizing Items

```starlark
# I'm more interested in Kubernetes SIG Auth issues than SIG API Machinery issues
# so I'd like to prioritize SIG Auth + SIG API Machinery issues as the highest, 
# SIG Auth as the next highest, then SIG API Machinery issues.
kubernetes_issues = github.search(query="repo:kubernetes/kubernetes is:open is:issue")

# When using the interactive output format, items
# are sorted by priority score in descending order
# (i.e highest priority first)
for issue in kubernetes_issues:
  prio = 0
  if "sig/auth" in issue.labels:
    prio += 20

  if "sig/api-machinery" in issue.labels:
    prio += 10

  issue.priority = prio

wranglr.render(kubernetes_issues)
```

## Multiple Sources

```starlark
kubernetes_issues = github.search(query="repo:kubernetes/kubernetes is:open is:issue")

openshift_bugs = jira.search(
  host="https://issues.redhat.com",
  query="project = \"OpenShift Bugs\""
)

wranglr.render(kubernetes_issues, openshift_bugs)
```
