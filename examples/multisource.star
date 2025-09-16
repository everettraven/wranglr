# Fetch all open issues and pull requests
# in the kubernetes-sigs/crdify repository
crdify_open_items = github.search(query="repo:kubernetes-sigs/crdify is:open")

# Fetch all jira tickets in the OpenShift Bugs
# project where kube-apiserver is included in
# the components
openshift_kubeapiserver_bugs = jira.search(
    host="https://issues.redhat.com",
    query="project = \"OpenShift Bugs\" AND component IN (kube-apiserver)",
)

# Tell wranglr to include both the kubernetes-sigs/crdify
# items and the OpenShift Bugs Jira tickets in the rendered output
wranglr.render(crdify_open_items, openshift_kubeapiserver_bugs)
