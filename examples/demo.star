sig_auth = github.search(query="repo:kubernetes/kubernetes is:open label:sig/auth", group="SIG Auth")
sig_apimachinery = github.search(query="repo:kubernetes/kubernetes is:open label:sig/api-machinery", group="SIG API Machinery")

kubeapilinter = github.search(query="repo:kubernetes-sigs/kube-api-linter is:open", group="API Tooling")
crdify = github.search(query="repo:kubernetes-sigs/crdify is:open", group="API Tooling")

all_issues = sig_auth + sig_apimachinery + kubeapilinter + crdify

for issue in all_issues:
  # pull requests vs issues
  if issue.pull_request["url"] != "":
    if "lgtm" not in issue.labels:
      issue.status = "Needs Review"
    elif "approved" not in issue.labels:
      issue.status = "Needs Approval"
  else:
    if "triage/accepted" not in issue.labels:
      issue.status = "Needs Triage"
    else:
      issue.status = "Todo"

wranglr.render(all_issues)
